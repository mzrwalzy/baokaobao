<template>
  <div class="dashboard">
    <div class="page-header">
      <h1>数据概览</h1>
    </div>

    <div class="stats-grid">
      <div class="stat-card">
        <div class="stat-icon" style="background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);">
          👥
        </div>
        <div class="stat-info">
          <div class="stat-value">{{ stats.total_users || 0 }}</div>
          <div class="stat-label">用户总数</div>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-icon" style="background: linear-gradient(135deg, #f59e0b 0%, #d97706 100%);">
          📝
        </div>
        <div class="stat-info">
          <div class="stat-value">{{ stats.total_questions || 0 }}</div>
          <div class="stat-label">题目总数</div>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-icon" style="background: linear-gradient(135deg, #10b981 0%, #059669 100%);">
          ✅
        </div>
        <div class="stat-info">
          <div class="stat-value">{{ stats.total_answers || 0 }}</div>
          <div class="stat-label">答题记录</div>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-icon" style="background: linear-gradient(135deg, #ec4899 0%, #db2777 100%);">
          🆕
        </div>
        <div class="stat-info">
          <div class="stat-value">{{ stats.today_users || 0 }}</div>
          <div class="stat-label">今日新增用户</div>
        </div>
      </div>
    </div>

    <div class="section-header">
      <h2>热门题库 TOP 5</h2>
    </div>

    <div class="bank-rank">
      <div v-for="(bank, index) in topBanks" :key="bank.bank_id" class="bank-rank-item">
        <div class="rank-badge" :class="`rank-${index + 1}`">{{ index + 1 }}</div>
        <div class="bank-info">
          <div class="bank-name">{{ bank.bank_name }}</div>
          <div class="bank-meta">
            <span>{{ bank.question_count }} 题</span>
            <span>{{ bank.purchased_count }} 人已购</span>
            <span>{{ bank.total_answers }} 次答题</span>
          </div>
        </div>
        <div class="bank-rate">
          <span :class="rateClass(bank.avg_correct_rate)">{{ (bank.avg_correct_rate || 0).toFixed(1) }}%</span>
          <div class="rate-label">平均正确率</div>
        </div>
      </div>
      <el-empty v-if="!topBanks.length" description="暂无数据" />
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getDashboard, getBankStatsList } from '@/api'

const stats = ref({})
const topBanks = ref([])

const loadData = async () => {
  try {
    const [dashRes, banksRes] = await Promise.all([
      getDashboard(),
      getBankStatsList({ page: 1, page_size: 5, sort_by: 'total_answers', sort_order: 'descending' })
    ])
    stats.value = dashRes
    topBanks.value = banksRes.list || []
  } catch (e) {
    console.error(e)
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
.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
  gap: 20px;
  margin-bottom: 32px;
}

.stat-card {
  background: #fff;
  border-radius: 12px;
  padding: 24px;
  display: flex;
  align-items: center;
  gap: 20px;
  transition: transform 0.2s, box-shadow 0.2s;

  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 8px 24px rgba(0,0,0,0.08);
  }

  .stat-icon {
    width: 60px;
    height: 60px;
    border-radius: 14px;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 28px;
  }

  .stat-info {
    .stat-value {
      font-size: 32px;
      font-weight: 700;
      color: #1f2937;
      line-height: 1.1;
    }

    .stat-label {
      font-size: 14px;
      color: #6b7280;
      margin-top: 4px;
    }
  }
}

.section-header {
  margin-bottom: 16px;

  h2 {
    font-size: 18px;
    font-weight: 600;
    color: #1f2937;
    margin: 0;
  }
}

.bank-rank {
  background: #fff;
  border-radius: 12px;
  padding: 20px;

  .bank-rank-item {
    display: flex;
    align-items: center;
    gap: 16px;
    padding: 14px 0;
    border-bottom: 1px solid #f3f4f6;

    &:last-child {
      border-bottom: none;
    }

    .rank-badge {
      width: 32px;
      height: 32px;
      border-radius: 8px;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 14px;
      font-weight: 700;
      color: #6b7280;
      background: #f3f4f6;
      flex-shrink: 0;

      &.rank-1 {
        background: #fef3c7;
        color: #d97706;
      }

      &.rank-2 {
        background: #e5e7eb;
        color: #4b5563;
      }

      &.rank-3 {
        background: #ffedd5;
        color: #c2410c;
      }
    }

    .bank-info {
      flex: 1;
      min-width: 0;

      .bank-name {
        font-size: 15px;
        font-weight: 600;
        color: #1f2937;
        margin-bottom: 4px;
      }

      .bank-meta {
        font-size: 13px;
        color: #6b7280;

        span {
          margin-right: 12px;

          &:last-child {
            margin-right: 0;
          }
        }
      }
    }

    .bank-rate {
      text-align: right;
      flex-shrink: 0;

      span {
        font-size: 18px;
        font-weight: 700;
      }

      .rate-label {
        font-size: 12px;
        color: #9ca3af;
        margin-top: 2px;
      }
    }
  }
}

.rate-high { color: #10b981; }
.rate-mid { color: #f59e0b; }
.rate-low { color: #ef4444; }
</style>
