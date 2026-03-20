<template>
  <div class="question-manage">
    <div class="page-header">
      <h1>题目管理</h1>
      <div class="header-actions">
        <el-button @click="handleImport">批量导入</el-button>
        <el-button type="primary" @click="handleCreate">新建题目</el-button>
      </div>
    </div>
    
    <div class="filters card">
      <el-form inline :model="filters">
        <el-form-item label="题库">
          <el-select v-model="filters.bank_id" placeholder="请选择" clearable>
            <el-option v-for="bank in banks" :key="bank.id" :label="bank.name" :value="bank.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="类型">
          <el-select v-model="filters.type" placeholder="请选择" clearable>
            <el-option label="单选题" value="single" />
            <el-option label="多选题" value="multiple" />
            <el-option label="判断题" value="truefalse" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>
    
    <div class="card">
      <el-table :data="list" v-loading="loading" stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="title" label="题目" min-width="250" show-overflow-tooltip />
        <el-table-column prop="type" label="类型" width="100">
          <template #default="{ row }">
            <el-tag size="small">{{ typeMap[row.type] || row.type }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="difficulty" label="难度" width="80">
          <template #default="{ row }">
            <el-rate v-model="row.difficulty" disabled size="small" />
          </template>
        </el-table-column>
        <el-table-column prop="bank_id" label="所属题库" width="120">
          <template #default="{ row }">
            {{ getBankName(row.bank_id) }}
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
              {{ row.status === 1 ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="handleEdit(row)">编辑</el-button>
            <el-button link type="primary" @click="handleView(row)">查看</el-button>
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
    
    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑题目' : '新建题目'" width="700px">
      <el-form :model="form" label-width="90px">
        <el-form-item label="所属题库" required>
          <el-select v-model="form.bank_id" placeholder="请选择">
            <el-option v-for="bank in banks" :key="bank.id" :label="bank.name" :value="bank.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="题目类型">
          <el-select v-model="form.type">
            <el-option label="单选题" value="single" />
            <el-option label="多选题" value="multiple" />
            <el-option label="判断题" value="truefalse" />
          </el-select>
        </el-form-item>
        <el-form-item label="题目标题" required>
          <el-input v-model="form.title" placeholder="请输入题目标题" />
        </el-form-item>
        <el-form-item label="题目内容">
          <el-input v-model="form.content" type="textarea" :rows="3" placeholder="题目描述或选项" />
        </el-form-item>
        <el-form-item label="正确答案" required>
          <el-input v-model="form.answer" placeholder="如: A 或 ABC" />
        </el-form-item>
        <el-form-item label="解析">
          <el-input v-model="form.analysis" type="textarea" :rows="2" placeholder="答案解析" />
        </el-form-item>
        <el-form-item label="难度">
          <el-rate v-model="form.difficulty" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>
    
    <el-dialog v-model="viewVisible" title="题目详情" width="600px">
      <div v-if="currentRow" class="question-detail">
        <p><strong>题目：</strong>{{ currentRow.title }}</p>
        <p><strong>内容：</strong>{{ currentRow.content }}</p>
        <p><strong>答案：</strong>{{ currentRow.answer }}</p>
        <p><strong>解析：</strong>{{ currentRow.analysis || '无' }}</p>
      </div>
    </el-dialog>
    
    <el-dialog v-model="importVisible" title="批量导入" width="600px">
      <div class="import-tips">
        <el-alert type="info" :closable="false">
          <template #title>
            导入说明：
            <ol style="margin: 8px 0 0 16px; padding-left: 0;">
              <li>请先下载 Excel 模板，按格式填写题目信息</li>
              <li>题库为必选项，请先选择要导入到的题库</li>
              <li>支持 .xlsx 格式的 Excel 文件</li>
            </ol>
          </template>
        </el-alert>
        <el-button type="primary" link @click="handleDownloadTemplate" style="margin-top: 12px;">
          📥 下载 Excel 模板
        </el-button>
      </div>
      <el-form style="margin-top: 20px;">
        <el-form-item label="题库" required>
          <el-select v-model="importForm.bank_id" placeholder="请选择题库">
            <el-option v-for="bank in banks" :key="bank.id" :label="bank.name" :value="bank.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="上传文件" required>
          <el-upload
            ref="uploadRef"
            :auto-upload="false"
            :limit="1"
            accept=".xlsx"
            :on-change="handleFileChange"
          >
            <template #trigger>
              <el-button>选择Excel文件</el-button>
            </template>
            <template #tip>
              <div class="el-upload__tip">支持 .xlsx 格式</div>
            </template>
          </el-upload>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="importVisible = false">取消</el-button>
        <el-button type="primary" :loading="importLoading" @click="handleImportSubmit">导入</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { getQuestions, createQuestion, updateQuestion, deleteQuestion, importQuestions, getQuestionBanks } from '@/api'
import { ElMessage, ElMessageBox } from 'element-plus'

const loading = ref(false)
const list = ref([])
const banks = ref([])
const page = ref(1)
const pageSize = ref(20)
const total = ref(0)
const dialogVisible = ref(false)
const viewVisible = ref(false)
const importVisible = ref(false)
const isEdit = ref(false)
const currentRow = ref(null)

const filters = reactive({ bank_id: null, type: '' })

const form = reactive({
  id: null,
  bank_id: null,
  title: '',
  content: '',
  answer: '',
  analysis: '',
  type: 'single',
  difficulty: 3
})

const importForm = reactive({
  bank_id: null,
  file: null
})

const importLoading = ref(false)

const typeMap = { single: '单选', multiple: '多选', truefalse: '判断' }

const loadData = async () => {
  loading.value = true
  try {
    const params = { page: page.value, page_size: pageSize.value }
    if (filters.bank_id) params.bank_id = filters.bank_id
    if (filters.type) params.type = filters.type
    const data = await getQuestions(params)
    list.value = data.list || []
    total.value = data.total || 0
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

const loadBanks = async () => {
  try {
    const data = await getQuestionBanks({ page: 1, page_size: 100 })
    banks.value = data.list || []
  } catch (e) {}
}

const getBankName = (id) => {
  const bank = banks.value.find(b => b.id === id)
  return bank ? bank.name : '-'
}

const handleSearch = () => { page.value = 1; loadData() }
const handleReset = () => { filters.bank_id = null; filters.type = ''; handleSearch() }

const handleCreate = () => {
  isEdit.value = false
  Object.assign(form, { id: null, bank_id: null, title: '', content: '', answer: '', analysis: '', type: 'single', difficulty: 3 })
  dialogVisible.value = true
}

const handleEdit = (row) => {
  isEdit.value = true
  Object.assign(form, row)
  dialogVisible.value = true
}

const handleView = (row) => {
  currentRow.value = row
  viewVisible.value = true
}

const handleSubmit = async () => {
  if (!form.title || !form.bank_id) {
    ElMessage.warning('请填写必填项')
    return
  }
  try {
    if (isEdit.value) {
      await updateQuestion(form.id, form)
      ElMessage.success('更新成功')
    } else {
      await createQuestion(form)
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
    await ElMessageBox.confirm(`确定删除题目"${row.title}"吗？`, '提示', { type: 'warning' })
    await deleteQuestion(row.id)
    ElMessage.success('删除成功')
    loadData()
  } catch (e) {
    if (e !== 'cancel') console.error(e)
  }
}

const handleImport = () => {
  importForm.bank_id = null
  importForm.file = null
  importVisible.value = true
}

const handleFileChange = (file) => {
  importForm.file = file.raw
}

const handleDownloadTemplate = () => {
  window.open('/admin/api/v1/question_template', '_blank')
}

const handleImportSubmit = async () => {
  if (!importForm.bank_id) {
    ElMessage.warning('请选择题库')
    return
  }
  if (!importForm.file) {
    ElMessage.warning('请上传文件')
    return
  }

  importLoading.value = true
  try {
    const formData = new FormData()
    formData.append('bank_id', importForm.bank_id)
    formData.append('file', importForm.file)

    await importQuestions(formData)
    ElMessage.success('导入成功')
    importVisible.value = false
    loadData()
  } catch (e) {
    ElMessage.error(e.message || '导入失败')
  } finally {
    importLoading.value = false
  }
}

onMounted(() => { loadData(); loadBanks() })
</script>

<style lang="scss" scoped>
.filters {
  margin-bottom: 16px;
  padding: 16px;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

.question-detail {
  p {
    margin: 12px 0;
    line-height: 1.6;
  }
}
</style>
