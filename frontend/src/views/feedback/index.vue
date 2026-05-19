<template>
  <div class="feedback-manage">
    <div class="page-header">
      <h1>用户反馈</h1>
    </div>

    <div class="card">
      <div class="filter-bar">
        <el-select
          v-model="filters.status"
          placeholder="处理状态"
          clearable
          style="width: 160px"
        >
          <el-option label="待处理" :value="0" />
          <el-option label="已处理" :value="1" />
        </el-select>
        <el-input
          v-model="filters.type"
          placeholder="反馈类型"
          clearable
          style="width: 180px"
          @keyup.enter="handleSearch"
        />
        <el-button type="primary" @click="handleSearch">查询</el-button>
        <el-button @click="handleReset">重置</el-button>
      </div>

      <el-table :data="list" v-loading="loading" stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="user_id" label="用户ID" width="100" />
        <el-table-column prop="type" label="反馈类型" width="140" />
        <el-table-column prop="content" label="反馈内容" min-width="200" show-overflow-tooltip />
        <el-table-column prop="contact" label="联系方式" width="140" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'warning'" size="small">
              {{ row.status === 1 ? '已处理' : '待处理' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="提交时间" width="180">
          <template #default="{ row }">
            {{ formatTime(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="120" fixed="right">
          <template #default="{ row }">
            <el-button
              link
              :type="row.status === 1 ? 'warning' : 'success'"
              @click="handleToggleStatus(row)"
            >
              {{ row.status === 1 ? '标记待处理' : '标记已处理' }}
            </el-button>
          </template>
        </el-table-column>
        <template #empty>
          <el-empty description="暂无反馈数据" />
        </template>
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
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { getFeedback, updateFeedbackStatus } from '@/api'
import { ElMessage } from 'element-plus'
import dayjs from 'dayjs'

const loading = ref(false)
const list = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(20)

const filters = reactive({
  status: null,
  type: ''
})

async function loadData() {
  loading.value = true
  try {
    const params = {
      page: page.value,
      page_size: pageSize.value
    }
    if (filters.status !== null && filters.status !== '') {
      params.status = filters.status
    }
    if (filters.type) params.type = filters.type
    const res = await getFeedback(params)
    list.value = res.list || []
    total.value = res.total || 0
  } catch {
    // error handled by request interceptor
  } finally {
    loading.value = false
  }
}

function handleSearch() {
  page.value = 1
  loadData()
}

function handleReset() {
  filters.status = null
  filters.type = ''
  page.value = 1
  loadData()
}

async function handleToggleStatus(row) {
  const newStatus = row.status === 1 ? 0 : 1
  const actionLabel = newStatus === 1 ? '已处理' : '待处理'
  try {
    await updateFeedbackStatus(row.id, { status: newStatus })
    ElMessage.success(`已标记为${actionLabel}`)
    row.status = newStatus
  } catch {
    // error handled by request interceptor
  }
}

function formatTime(time) {
  return time ? dayjs(time).format('YYYY-MM-DD HH:mm') : '-'
}

onMounted(() => {
  loadData()
})
</script>

<style lang="scss" scoped>
.feedback-manage {
  .page-header {
    margin-bottom: 24px;

    h1 {
      font-size: 20px;
      font-weight: 600;
      color: #1f2937;
      margin: 0;
    }
  }

  .card {
    background: #fff;
    border-radius: 8px;
    padding: 24px;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.06);
  }

  .filter-bar {
    display: flex;
    align-items: center;
    gap: 12px;
    margin-bottom: 20px;
    flex-wrap: wrap;
  }

  .pagination {
    display: flex;
    justify-content: flex-end;
    margin-top: 20px;
  }
}
</style>
