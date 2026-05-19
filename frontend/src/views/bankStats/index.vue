<template>
  <div class="bank-stats">
    <div class="page-header">
      <h1>题库表现统计</h1>
    </div>

    <div class="card">
      <el-table
        :data="list"
        v-loading="loading"
        stripe
        :default-sort="{ prop: 'avg_correct_rate', order: 'ascending' }"
        @sort-change="handleSortChange"
        @row-click="handleRowClick"
        highlight-current-row
      >
        <el-table-column prop="bank_name" label="题库名称" min-width="160" />
        <el-table-column prop="question_count" label="题目数" width="100" sortable="custom" />
        <el-table-column prop="purchased_count" label="已购人数" width="100" sortable="custom" />
        <el-table-column prop="total_answers" label="答题次数" width="110" sortable="custom" />
        <el-table-column prop="total_exams" label="考试次数" width="110" sortable="custom" />
        <el-table-column prop="avg_correct_rate" label="平均正确率" width="130" sortable="custom">
          <template #default="{ row }">
            <span :class="rateClass(row.avg_correct_rate)">
              {{ (row.avg_correct_rate || 0).toFixed(1) }}%
            </span>
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

    <el-drawer
      v-model="drawerVisible"
      :title="`${currentBankName || '题库'} - 错题详情`"
      size="560px"
    >
      <div v-if="bankStats" v-loading="statsLoading" class="bank-stats-detail">
        <div class="stats-grid">
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
        <div v-if="bankStats.top_wrong_questions?.length" class="wrong-questions">
          <h5>错误率最高题目 TOP10</h5>
          <el-table :data="bankStats.top_wrong_questions" size="small" stripe>
            <el-table-column type="index" label="排名" width="60" />
            <el-table-column prop="title" label="题目标题" min-width="160" show-overflow-tooltip />
            <el-table-column prop="content" label="内容" min-width="200" show-overflow-tooltip />
            <el-table-column prop="wrong_count" label="错误次数" width="90" />
          </el-table>
        </div>
        <el-empty v-else description="暂无错题数据" />
      </div>
    </el-drawer>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getBankStatsList, getBankStats } from '@/api'

const loading = ref(false)
const list = ref([])
const page = ref(1)
const pageSize = ref(20)
const total = ref(0)
const sortProp = ref('')
const sortOrder = ref('')

const drawerVisible = ref(false)
const currentBankId = ref(null)
const currentBankName = ref('')
const bankStats = ref(null)
const statsLoading = ref(false)

const loadData = async () => {
  loading.value = true
  try {
    const params = { page: page.value, page_size: pageSize.value }
    if (sortProp.value) {
      params.sort_by = sortProp.value
      params.sort_order = sortOrder.value
    }
    const data = await getBankStatsList(params)
    list.value = data.list || []
    total.value = data.total || 0
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

const handleSortChange = ({ prop, order }) => {
  sortProp.value = prop || ''
  sortOrder.value = order || ''
  page.value = 1
  loadData()
}

const handleRowClick = async (row) => {
  currentBankId.value = row.bank_id
  currentBankName.value = row.bank_name
  drawerVisible.value = true
  await loadBankStats(row.bank_id)
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

const rateClass = (rate) => {
  if (rate >= 80) return 'rate-high'
  if (rate >= 60) return 'rate-mid'
  return 'rate-low'
}

onMounted(loadData)
</script>

<style lang="scss" scoped>
.bank-stats {
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
}

.bank-stats-detail {
  .stats-grid {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 12px;
    margin-bottom: 20px;

    .stat-item {
      background: #f9fafb;
      border-radius: 8px;
      padding: 12px;
      text-align: center;

      .stat-value {
        font-size: 20px;
        font-weight: 700;
        color: #1f2937;
      }

      .stat-label {
        font-size: 12px;
        color: #6b7280;
        margin-top: 4px;
      }
    }
  }

  .wrong-questions {
    margin-top: 16px;

    h5 {
      font-size: 13px;
      font-weight: 600;
      color: #4b5563;
      margin-bottom: 8px;
    }
  }
}
</style>
