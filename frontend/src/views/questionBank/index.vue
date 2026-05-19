<template>
  <div class="question-bank">
    <div class="page-header">
      <h1>题库管理</h1>
      <el-button type="primary" @click="handleCreate">新建题库</el-button>
    </div>

    <div class="card">
      <el-table :data="list" v-loading="loading" stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="题库名称" min-width="150" />
        <el-table-column prop="description" label="描述" min-width="200" show-overflow-tooltip />
        <el-table-column prop="price" label="价格" width="100">
          <template #default="{ row }">
            ¥{{ row.price || 0 }}
          </template>
        </el-table-column>
        <el-table-column prop="question_num" label="题目数" width="90" />
        <el-table-column prop="purchased_count" label="已购人数" width="90">
          <template #default="{ row }">
            {{ row.purchased_count !== undefined ? row.purchased_count : '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="total_answers" label="答题次数" width="90">
          <template #default="{ row }">
            {{ row.total_answers !== undefined ? row.total_answers : '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="avg_correct_rate" label="平均正确率" width="110">
          <template #default="{ row }">
            <span :class="rateClass(row.avg_correct_rate)">
              {{ row.avg_correct_rate !== undefined ? (row.avg_correct_rate).toFixed(1) + '%' : '-' }}
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="90">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
              {{ row.status === 1 ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="170">
          <template #default="{ row }">
            {{ formatTime(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="handleEdit(row)">编辑</el-button>
            <el-button link type="info" @click="handleStats(row)">统计</el-button>
            <el-button link type="danger" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination">
        <el-pagination
          v-model:current-page="page"
          v-model:page-size="pageSize"
          :total="total"
          :page-sizes="[10, 20, 50]"
          layout="total, sizes, prev, pager, next"
          @change="loadData"
        />
      </div>
    </div>
    
    <el-dialog
      v-model="dialogVisible"
      :title="isEdit ? '编辑题库' : '新建题库'"
      width="500px"
    >
      <el-form :model="form" label-width="80px">
        <el-form-item label="题库名称" required>
          <el-input v-model="form.name" placeholder="请输入题库名称" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="form.description" type="textarea" :rows="3" placeholder="请输入描述" />
        </el-form-item>
        <el-form-item label="价格">
          <el-input-number v-model="form.price" :min="0" :precision="2" />
        </el-form-item>
        <el-form-item label="封面图">
          <el-input v-model="form.cover_image" placeholder="封面图片URL" />
        </el-form-item>
        <el-form-item label="状态">
          <el-radio-group v-model="form.status">
            <el-radio :label="1">启用</el-radio>
            <el-radio :label="0">禁用</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>

    <el-drawer
      v-model="drawerVisible"
      :title="`${currentBankName || '题库'} - 统计详情`"
      size="560px"
    >
      <div v-if="bankStats" v-loading="statsLoading" class="bank-stats-detail">
        <div class="stats-grid">
          <div class="stat-item">
            <div class="stat-value">{{ bankStats.total_users || 0 }}</div>
            <div class="stat-label">已购人数</div>
          </div>
          <div class="stat-item">
            <div class="stat-value">{{ bankStats.total_answers || 0 }}</div>
            <div class="stat-label">答题次数</div>
          </div>
          <div class="stat-item">
            <div class="stat-value">{{ bankStats.total_exams || 0 }}</div>
            <div class="stat-label">考试次数</div>
          </div>
          <div class="stat-item">
            <div class="stat-value">{{ (bankStats.avg_correct_rate || 0).toFixed(1) }}%</div>
            <div class="stat-label">平均正确率</div>
          </div>
        </div>
        <div v-if="bankStats.difficulty_dist?.length" class="dist-section">
          <h5>难度分布</h5>
          <div class="dist-list">
            <div v-for="item in bankStats.difficulty_dist" :key="item.difficulty" class="dist-row">
              <span class="dist-label">难度 {{ item.difficulty }}</span>
              <el-progress :percentage="(item.count / bankStats.total_questions * 100).toFixed(0)" :format="() => item.count + '题'" :stroke-width="8" />
            </div>
          </div>
        </div>
        <div v-if="bankStats.daily_trend?.length" class="trend-section">
          <h5>近30日趋势</h5>
          <div class="trend-list">
            <div v-for="item in bankStats.daily_trend.slice(-7)" :key="item.date" class="trend-row">
              <span class="trend-date">{{ item.date }}</span>
              <span class="trend-count">{{ item.answer_count }}次</span>
              <span :class="['trend-rate', rateClass(item.correct_rate)]">{{ item.correct_rate.toFixed(1) }}%</span>
            </div>
          </div>
        </div>
        <el-empty v-else description="暂无趋势数据" />
      </div>
    </el-drawer>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { getQuestionBanks, createQuestionBank, updateQuestionBank, deleteQuestionBank, getBankStatsList, getBankStats } from '@/api'
import { ElMessage, ElMessageBox } from 'element-plus'
import dayjs from 'dayjs'

const loading = ref(false)
const list = ref([])
const page = ref(1)
const pageSize = ref(20)
const total = ref(0)
const dialogVisible = ref(false)
const isEdit = ref(false)

const drawerVisible = ref(false)
const currentBankId = ref(null)
const currentBankName = ref('')
const bankStats = ref(null)
const statsLoading = ref(false)

const form = reactive({
  id: null,
  name: '',
  description: '',
  price: 0,
  cover_image: '',
  status: 1
})

const loadData = async () => {
  loading.value = true
  try {
    const [banksRes, statsRes] = await Promise.all([
      getQuestionBanks({ page: page.value, page_size: pageSize.value }),
      getBankStatsList({ page: 1, page_size: 1000 }) // 拉取全部用于合并
    ])

    const banks = banksRes.list || []
    const statsList = statsRes.list || []
    const statsMap = new Map(statsList.map(s => [s.bank_id, s]))

    list.value = banks.map(bank => {
      const stat = statsMap.get(bank.id)
      return {
        ...bank,
        purchased_count: stat?.purchased_count,
        total_answers: stat?.total_answers,
        avg_correct_rate: stat?.avg_correct_rate
      }
    })
    total.value = banksRes.total || 0
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

const handleCreate = () => {
  isEdit.value = false
  Object.assign(form, { id: null, name: '', description: '', price: 0, cover_image: '', status: 1 })
  dialogVisible.value = true
}

const handleEdit = (row) => {
  isEdit.value = true
  Object.assign(form, row)
  dialogVisible.value = true
}

const handleStats = async (row) => {
  currentBankId.value = row.id
  currentBankName.value = row.name
  drawerVisible.value = true
  await loadBankStats(row.id)
}

const loadBankStats = async (bankId) => {
  statsLoading.value = true
  try {
    const data = await getBankStats(bankId)
    bankStats.value = data
  } catch (e) {
    console.error(e)
  } finally {
    statsLoading.value = false
  }
}

const handleSubmit = async () => {
  if (!form.name) {
    ElMessage.warning('请输入题库名称')
    return
  }
  try {
    if (isEdit.value) {
      await updateQuestionBank(form.id, form)
      ElMessage.success('更新成功')
    } else {
      await createQuestionBank(form)
      ElMessage.success('创建成功')
    }
    dialogVisible.value = false
    loadData()
  } catch (e) {
    console.error(e)
  }
}

const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm(`确定删除题库"${row.name}"吗？`, '提示', { type: 'warning' })
    await deleteQuestionBank(row.id)
    ElMessage.success('删除成功')
    loadData()
  } catch (e) {
    if (e !== 'cancel') console.error(e)
  }
}

const rateClass = (rate) => {
  if (rate >= 80) return 'rate-high'
  if (rate >= 60) return 'rate-mid'
  return 'rate-low'
}

const formatTime = (time) => time ? dayjs(time).format('YYYY-MM-DD HH:mm') : '-'

onMounted(loadData)
</script>

<style lang="scss" scoped>
.card {
  background: #fff;
  border-radius: 12px;
  padding: 20px;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

.rate-high {
  color: #10b981;
  font-weight: 600;
}

.rate-mid {
  color: #f59e0b;
  font-weight: 600;
}

.rate-low {
  color: #ef4444;
  font-weight: 600;
}

.bank-stats-detail {
  .stats-grid {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: 12px;
    margin-bottom: 24px;

    .stat-item {
      background: #f9fafb;
      border-radius: 8px;
      padding: 16px;
      text-align: center;

      .stat-value {
        font-size: 24px;
        font-weight: 700;
        color: #1f2937;
      }

      .stat-label {
        font-size: 13px;
        color: #6b7280;
        margin-top: 4px;
      }
    }
  }

  .dist-section,
  .trend-section {
    margin-top: 20px;

    h5 {
      font-size: 14px;
      font-weight: 600;
      color: #374151;
      margin-bottom: 12px;
    }
  }

  .dist-list {
    .dist-row {
      display: flex;
      align-items: center;
      gap: 12px;
      margin-bottom: 8px;

      .dist-label {
        width: 60px;
        font-size: 13px;
        color: #6b7280;
        flex-shrink: 0;
      }
    }
  }

  .trend-list {
    .trend-row {
      display: flex;
      align-items: center;
      justify-content: space-between;
      padding: 8px 0;
      border-bottom: 1px solid #f3f4f6;

      .trend-date {
        font-size: 13px;
        color: #6b7280;
      }

      .trend-count {
        font-size: 13px;
        color: #374151;
      }

      .trend-rate {
        font-size: 13px;
        font-weight: 600;
      }
    }
  }
}
</style>
