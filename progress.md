# Progress Log

## 2026-05-14

### Completed
- [x] Read and analyzed `docs/admin-management-roadmap.md`
- [x] Created planning files (`task_plan.md`, `findings.md`, `progress.md`)
- [x] Phase 3.2 + 3.3: User study records & exam records views (backend APIs + frontend drawer tabs)
- [x] Phase 4.1: Import result details (backend returns success/failure/failures + frontend result dialog)
- [x] Phase 2.2: Bank deletion strategy (backend dependency check + frontend fallback to disable)
- [x] Phase 2.3: Bank operations info (purchased_count in list + detail drawer)
- [x] Phase 5.1: Dashboard core metrics (new stat cards + correct rate + 7-day trends)
- [x] Phase 5.2: Bank performance stats (per-bank API + detail drawer stats + independent `/bank-stats` page)
- [x] Phase 6.1: Unified form validation (frontend rules + backend field validation)
- [x] OSS integration (Aliyun OSS SDK + upload logic with local fallback)
- [x] Phase 4.3: Question export (backend Excel export + frontend button)
- [x] Phase 6.2: Role-based access control (superadmin vs operator, RoleMiddleware, route protection, frontend role store)
- [x] Phase 6.3: Operation audit logs (AuditLog model + audit_logs table + recording on key actions + `/audit-log` frontend page with filters)
- [x] Verified builds: Go backend compiles cleanly, Vue frontend builds successfully

### Files Modified (summary)
- **Backend**: `internal/handler/handler.go`, `internal/router/router.go`, `internal/model/*.go`, `internal/repository/*.go`, `internal/service/*.go`, `internal/middleware/role.go`, `internal/pkg/oss/oss.go`
- **Frontend**: `frontend/src/views/*` (user, questionBank, question, dashboard, bankStats, auditLog), `frontend/src/api/index.js`, `frontend/src/router/index.js`, `frontend/src/layouts/BasicLayout.vue`, `frontend/src/stores/user.js`
- **Config**: `go.mod` (added aliyun-oss-go-sdk), `docs/admin-management-roadmap.md`

### Next
后台 roadmap 全部完成。后续可选：
1. 前端按角色动态隐藏按钮/菜单
2. 题目删除增加依赖检查
3. 小程序端功能联调
4. OSS 生产环境配置
5. 统计查询性能优化（索引/缓存）
