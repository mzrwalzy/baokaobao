# 代码仓库审查报告

- 审查日期: 2026-05-19
- 审查基线: `c67dd0c chore: snapshot current changes before audit`
- 审查范围: Go 后端、Vue 管理端、微信小程序、仓库工程化文件
- 已执行检查:
  - `go test ./...`
  - `go vet ./...`
  - `npm run build` in `frontend/`
  - 关键路径人工审查: 鉴权、权限、统计、审计日志、题库访问、小程序请求层

## 总结

本次审查没有发现会直接阻止后端编译或前端构建的问题，但发现了多处高风险的权限与契约缺陷。最关键的问题集中在三类:

1. 后台角色权限没有真正落地，`operator` 与 `superadmin` 的边界形同虚设。
2. JWT 会话只在登录时校验状态，账号被禁用后，已签发 token 仍可继续访问。
3. 若干前后端新增功能只完成了一半，接口和页面字段不一致，或能力已经存表但从未真正参与授权。

## 发现列表

### P0: 后台 RBAC 未生效，任意管理员可访问所有后台能力

- 严重级别: `Critical`
- 影响:
  - `operator` 可以直接调用管理员管理、密码重置、题库权限发放、审计日志、反馈处理等所有接口。
  - 角色隔离只停留在 token 字段和前端本地状态，后端没有真正执行授权边界。
- 证据:
  - `RoleMiddleware` 已实现，但未挂载到任何路由: `internal/middleware/role.go:9-11`
  - 后台路由统一只用了 `AdminAuth()`，写操作组也没有角色限制: `internal/router/router.go:128-165`
  - 前端路由守卫只检查是否有 token，不检查角色: `frontend/src/router/index.js:67-75`
- 风险说明:
  - 这是典型的越权问题，属于后台安全边界失效。
  - `progress.md` 已记录“Role-based access control 完成”，但当前实现与记录不一致: `progress.md:17`

### P0: 被禁用用户/管理员的已签发 token 不会立即失效

- 严重级别: `Critical`
- 影响:
  - 小程序用户和后台管理员都只在登录时检查 `status`。
  - 账号被禁用后，只要旧 token 未过期，仍可继续访问受保护接口。
- 证据:
  - 登录时检查状态: `internal/service/auth_service.go:68-69`, `internal/service/auth_service.go:104-106`
  - 中间件只校验 JWT 类型与黑名单，不回表检查当前账号状态: `internal/middleware/jwt_auth.go:25-65`, `internal/middleware/jwt_auth.go:69-110`
- 风险说明:
  - 管理后台场景下，这会让“封禁账号”“停用离职管理员”失去即时性。

### P1: 管理员题库权限功能是“死功能”，数据写入后并未参与任何授权

- 严重级别: `High`
- 影响:
  - 系统可以创建、查询、删除 `admin_bank_permissions` 记录，但这些记录从未参与题库列表、题目列表、编辑、导入、删除等实际权限判断。
  - 最终效果是“有权限表、无权限控制”。
- 证据:
  - 仅有 CRUD 能力: `internal/handler/admin_handler.go:178-246`, `internal/service/admin_service.go:145-154`, `internal/repository/admin_repo.go:144-159`
  - 后台题库/题目相关路由没有任何按角色或按题库的细粒度授权: `internal/router/router.go:131-165`
  - 全仓检索未发现这些权限数据被题库/题目查询与写入逻辑消费。
- 风险说明:
  - 功能表面可用，实际不生效，容易误导运营人员。

### P1: 小程序生产地址硬编码为本地 `localhost`

- 严重级别: `High`
- 影响:
  - 小程序一旦切到真实环境，所有 API 调用都会指向本地开发地址，接口整体不可用。
- 证据:
  - `miniprogram/utils/request.js:1-2`
- 风险说明:
  - 这是明确的发布阻塞项。

### P1: 审计日志页面与后端接口字段不一致，列表展示与筛选都不完整

- 严重级别: `High`
- 影响:
  - 前端展示列读取 `resource` / `resource_id`，而后端记录的是 `target` / `target_id`，对应列会为空或回退原值。
  - 前端传 `start_date` / `end_date`，后端读取 `start_time` / `end_time`，日期筛选无效。
- 证据:
  - 前端列与筛选参数: `frontend/src/views/auditLog/index.vue:51-56`, `frontend/src/views/auditLog/index.vue:153-154`
  - 后端记录字段与查询参数: `internal/handler/handler.go:990-1004`, `internal/handler/handler.go:1023-1031`
- 风险说明:
  - 审计日志本来就是安全与追责工具，这类字段错位会直接降低可用性。

### P1: 小程序考试记录分页参数传错，翻页逻辑实际上拿不到下一页

- 严重级别: `High`
- 影响:
  - 页面翻页时调用 `getExamRecords({ page, page_size })`，但 API 包装成了 `{ params }`。
  - `wx.request` 层不会像 axios 那样解包 `params`，后端拿不到 `page` / `page_size`，分页会退回默认值。
- 证据:
  - 小程序 API 包装: `miniprogram/api/index.js:19`
  - 请求层直接透传 `data`: `miniprogram/utils/request.js:9-12`, `miniprogram/utils/request.js:52`
  - 页面依赖服务端分页结果更新游标: `miniprogram/pages/profile/exam-records.js:54-71`
- 风险说明:
  - 表现上通常是“滚动加载看起来工作，但永远在重复第一页数据”。

### P2: 访问授权表缺少唯一约束，重复授权会污染统计与权限状态

- 严重级别: `Medium`
- 影响:
  - `user_bank_access` 与 `admin_bank_permissions` 都没有联合唯一键。
  - 仓储层也直接 `Create`，没有幂等保护。
  - 一旦重复插入，会导致已购人数、题库用户数、权限状态判断等结果失真。
- 证据:
  - 模型缺少联合唯一索引: `internal/model/model.go:154-160`, `internal/model/admin_bank_permission.go:5-9`
  - 仓储层直接写入: `internal/repository/user_repo.go:46-47`, `internal/repository/admin_repo.go:150-155`
  - 统计查询直接按表计数: `internal/repository/bank_stats_repo.go:103-104`
- 风险说明:
  - 这是典型的数据一致性问题，越往后修复成本越高。

### P2: 后台题库统计独立页没有真正接入导航，且页面预期字段大于后端实际返回

- 严重级别: `Medium`
- 影响:
  - 仓库里存在完整的 `bankStats` 页面，但路由并未注册，也没有菜单入口。
  - 页面详情抽屉还依赖 `top_wrong_questions`，而后端 `BankDetailStat` 未返回该字段。
- 证据:
  - 页面存在并读取 `top_wrong_questions`: `frontend/src/views/bankStats/index.vue:63-70`
  - 路由未注册该页面: `frontend/src/router/index.js:4-59`
  - 后端返回结构不含 `top_wrong_questions`: `internal/repository/bank_stats_repo.go:17-27`
- 风险说明:
  - 这属于“功能完成度被高估”，容易在验收或联调阶段暴露。

### P2: `UpdateAdminUser` 的部分更新语义不安全，存在把角色清空的可能

- 严重级别: `Medium`
- 影响:
  - Handler 允许 `role` 为空字符串。
  - Repository 无条件把 `role` 写回数据库。
  - 只做昵称/状态更新的第三方调用若未带 `role`，可能把管理员角色清空。
- 证据:
  - Handler 允许空角色通过: `internal/handler/admin_handler.go:95-115`
  - Repository 无条件更新 `role`: `internal/repository/admin_repo.go:128-133`
- 风险说明:
  - 当前前端页面会传 `role`，但 API 本身不安全，后续接入脚本/别的前端时容易出错。

### P3: 仓库已提交调试/草稿文件，影响代码整洁度与后续维护

- 严重级别: `Low`
- 影响:
  - `internal/handler/handler.go.my-edits` 已被 Git 跟踪，容易让后续开发误判哪份代码才是有效实现。
- 证据:
  - `git ls-files` 结果包含 `internal/handler/handler.go.my-edits`
- 风险说明:
  - 虽然不影响编译，但属于典型的仓库卫生问题，应尽快清理并加入忽略规则。

## 未列为缺陷但建议关注

- 前端构建通过，但 `vite build` 提示存在大体积 chunk，后续可列入性能治理。
- 小程序里仍有不少页面只做 `console.error` 而不提示用户，体验层面还有统一错误处理空间。
- 当前统计查询以联表聚合为主，随着数据增长可能需要索引与缓存治理。
