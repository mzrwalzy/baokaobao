# Findings Log

## Session Start
- Date: 2026-05-14
- Source: `docs/admin-management-roadmap.md`
- Priority tasks: 3.2 (Study Records), 3.3 (Exam Records), 4.1 (Import Details)

## Codebase Structure (from README)
- Frontend: Vue 3 + Vite + Element Plus, Composition API with `<script setup>`
- Backend: Go 1.21+, Gin, GORM, MySQL
- Admin API base: `/admin/api/v1/`
- Existing user answers table: `user_answers`
- Existing exam records table: `exam_records`
- Existing wrong questions table: `wrong_questions`

## To Investigate
- How are user details currently displayed in admin?
- What quiz/exam record structures exist in backend?
- How does Excel import currently work?

## Notes
- Follow existing `handler → service → repository` pattern
- Use `response.Page()` for paginated admin APIs
- Frontend uses `frontend/src/api/index.js` for API calls
- Response interceptor expects `res.code === 0`
