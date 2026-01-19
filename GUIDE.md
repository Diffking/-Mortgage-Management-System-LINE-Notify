# 📘 SPSC loanEasy v1.0

> ระบบสินเชื่อสหกรณ์ออมทรัพย์สาธารณสุขจังหวัดสงขลา

```
╔═══════════════════════════════════════════════════════════════╗
║                    🏦 SPSC loanEasy API                       ║
║                        Version 1.0                            ║
╠═══════════════════════════════════════════════════════════════╣
║  ✅ Phase 1: Project Setup + Database + Swagger               ║
║  ✅ Phase 2: Authentication (JWT + Refresh Token)             ║
║  ✅ Phase 3: User Management + RBAC                           ║
║  ✅ Phase 4: Mortgage Management + LINE Notify                ║
║  ⏳ Phase 5: LINE OA แจ้งเตือนรายบุคคล                         ║
╚═══════════════════════════════════════════════════════════════╝
```

---

## 📋 สารบัญ

| หัวข้อ | คำอธิบาย |
|--------|----------|
| [🔧 Prerequisites](#-prerequisites) | ซอฟต์แวร์ที่ต้องติดตั้ง |
| [📥 Installation](#-installation) | ขั้นตอนการติดตั้ง |
| [⚙️ Configuration](#️-configuration) | การตั้งค่า Environment |
| [🚀 Running](#-running) | การรันโปรเจกต์ |
| [🔐 Authentication](#-authentication-phase-2) | ระบบ Login/Register |
| [👥 User Management](#-user-management-phase-3) | จัดการผู้ใช้ + RBAC |
| [🏦 Mortgage](#-mortgage-phase-4) | ระบบจำนอง |
| [📡 API Reference](#-api-reference) | รายการ API ทั้งหมด |
| [🔍 Troubleshooting](#-troubleshooting) | แก้ไขปัญหา |

---

## 🔧 Prerequisites

```
┌─────────────────────────────────────────────────────────────┐
│  💻 Required Software                                       │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│   Go          ≥ 1.21      →  go version                     │
│   MySQL       ≥ 8.0       →  mysql --version                │
│   Git         Latest      →  git --version                  │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

---

## 📥 Installation

```bash
# 1. แตกไฟล์
unzip spsc-loaneasy-phase4.zip -d spsc-loaneasy
cd spsc-loaneasy

# 2. สร้าง Database
mysql -u root -p -e "CREATE DATABASE loaneasy CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"

# 3. ตั้งค่า .env
cp .env.example .env
# แก้ไขค่า DEV_DB_PASS

# 4. ติดตั้ง Dependencies
go mod tidy

# 5. รัน Server
set APP_MODE=dev && go run cmd/server/main.go
```

---

## ⚙️ Configuration

แก้ไขไฟล์ `.env`:

```env
# Application
APP_MODE=dev
PORT=3000

# Database
DEV_DB_HOST=localhost
DEV_DB_PORT=3306
DEV_DB_USER=root
DEV_DB_PASS=your_password
DEV_DB_NAME=loaneasy

# JWT
DEV_JWT_SECRET=your_secret_key_min_32_chars
DEV_JWT_REFRESH_SECRET=your_refresh_secret_key

# Cookie
DEV_COOKIE_SECURE=false

# LINE Notify (Phase 4)
LINE_NOTIFY_TOKEN=your_line_notify_token
```

---

## 🏦 Mortgage (Phase 4)

### 📊 Database Schema

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                            📊 ERD - Phase 4                                 │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   🔵 Master Tables (4)                                                      │
│   ├── loan_types      : ประเภทเงินกู้ (สามัญ, ฉุกเฉิน, พิเศษ)               │
│   ├── loan_steps      : ขั้นตอน (ร่าง→รอเอกสาร→ตรวจสอบ→อนุมัติ/ปฏิเสธ)      │
│   ├── loan_docs       : รายการเอกสาร (บัตรปชช, ทะเบียนบ้าน, ฯลฯ)           │
│   └── loan_appts      : ประเภทนัดหมาย (เซ็นสัญญา, ส่งเอกสาร, รับเงิน)       │
│                                                                             │
│   🟢 Main Tables (2)                                                        │
│   ├── mortgages       : ข้อมูลจำนอง (สมาชิกดูได้, เจ้าหน้าที่ Create)       │
│   └── transactions    : History ทุกการเปลี่ยนแปลง                          │
│                                                                             │
│   🟡 Current Tables (4)                                                     │
│   ├── loan_doc_currents   : Checklist เอกสาร (1:N)                         │
│   ├── loan_type_currents  : ประเภท ณ เวลานั้น                               │
│   ├── loan_step_currents  : สถานะ ณ เวลานั้น                                │
│   └── loan_appt_currents  : นัดหมาย (วันที่จริง)                            │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 👥 Roles & Permissions

```
┌────────────────────────────────────────────────────────────────┐
│                    👥 Role-Based Access Control                │
├────────────────────────────────────────────────────────────────┤
│                                                                │
│   🟢 USER (สมาชิก)                                             │
│   └── ดูจำนองของตัวเอง (GET /mortgages/my)                    │
│                                                                │
│   🟡 OFFICER (เจ้าหน้าที่)                                     │
│   └── ทุกอย่างของ USER                                         │
│   └── สร้าง/แก้ไข mortgage                                     │
│   └── เปลี่ยนสถานะ, ตรวจเอกสาร, นัดหมาย                        │
│   └── อนุมัติ/ปฏิเสธ                                           │
│                                                                │
│   🔴 ADMIN                                                     │
│   └── ทุกอย่างของ OFFICER                                      │
│   └── จัดการ Master Data                                       │
│   └── เปลี่ยนเจ้าหน้าที่ผู้รับผิดชอบ                            │
│                                                                │
└────────────────────────────────────────────────────────────────┘
```

### 🔔 LINE Notify

```
┌─────────────────────────────────────────────────────────────┐
│  แจ้งเตือนอัตโนมัติเมื่อ:                                    │
│                                                             │
│  🆕 สร้างคำขอใหม่                                           │
│  🔄 เปลี่ยนสถานะ                                            │
│  ✅ อนุมัติ                                                 │
│  ❌ ปฏิเสธ                                                  │
│  📅 สร้างนัดหมาย                                            │
│                                                             │
│  ตั้งค่า: LINE_NOTIFY_TOKEN ใน .env                         │
│  สร้าง Token: https://notify-bot.line.me/                   │
└─────────────────────────────────────────────────────────────┘
```

---

## 📡 API Reference

### 🔐 Authentication

| Method | Endpoint | Auth | Description |
|:------:|----------|:----:|-------------|
| `POST` | `/api/v1/auth/register` | ❌ | ลงทะเบียน |
| `POST` | `/api/v1/auth/login` | ❌ | เข้าสู่ระบบ |
| `POST` | `/api/v1/auth/refresh` | ❌ | รีเฟรช Token |
| `POST` | `/api/v1/auth/logout` | ❌ | ออกจากระบบ |
| `GET` | `/api/v1/auth/me` | ✅ | ดูข้อมูลตัวเอง |

### 👤 Profile

| Method | Endpoint | Auth | Description |
|:------:|----------|:----:|-------------|
| `GET` | `/api/v1/profile` | ✅ | ดูโปรไฟล์ |
| `PUT` | `/api/v1/profile` | ✅ | แก้ไขโปรไฟล์ |
| `PUT` | `/api/v1/profile/password` | ✅ | เปลี่ยนรหัสผ่าน |

### 👥 Users (Admin)

| Method | Endpoint | Auth | Description |
|:------:|----------|:----:|-------------|
| `GET` | `/api/v1/users` | 🔴 | รายการผู้ใช้ |
| `GET` | `/api/v1/users/:id` | 🔴 | ดูผู้ใช้ |
| `PUT` | `/api/v1/users/:id` | 🔴 | แก้ไขผู้ใช้ |
| `DELETE` | `/api/v1/users/:id` | 🔴 | ลบผู้ใช้ |
| `PUT` | `/api/v1/users/:id/role` | 🔴 | เปลี่ยน Role |

### 🏦 Mortgages

| Method | Endpoint | Auth | Description |
|:------:|----------|:----:|-------------|
| `GET` | `/api/v1/mortgages/my` | ✅ | ดูจำนองของตัวเอง |
| `POST` | `/api/v1/mortgages` | 🟡 | สร้างจำนอง |
| `GET` | `/api/v1/mortgages` | 🟡 | รายการจำนอง |
| `GET` | `/api/v1/mortgages/:id` | 🟡 | ดูรายละเอียด |
| `GET` | `/api/v1/mortgages/:id/history` | 🟡 | ดู History |
| `PUT` | `/api/v1/mortgages/:id/step` | 🟡 | เปลี่ยนสถานะ |
| `PUT` | `/api/v1/mortgages/:id/approve` | 🟡 | อนุมัติ |
| `PUT` | `/api/v1/mortgages/:id/reject` | 🟡 | ปฏิเสธ |
| `GET` | `/api/v1/mortgages/:id/docs` | 🟡 | ดูเอกสาร |
| `PUT` | `/api/v1/mortgages/:id/docs` | 🟡 | อัพเดทเอกสาร |
| `GET` | `/api/v1/mortgages/:id/appts` | 🟡 | ดูนัดหมาย |
| `POST` | `/api/v1/mortgages/:id/appts` | 🟡 | สร้างนัดหมาย |
| `PUT` | `/api/v1/mortgages/:id/appts/:appt_id/complete` | 🟡 | นัดหมายเสร็จ |
| `PUT` | `/api/v1/mortgages/:id/officer` | 🔴 | เปลี่ยนเจ้าหน้าที่ |

### 📋 Master Data (Admin)

| Method | Endpoint | Description |
|:------:|----------|-------------|
| `CRUD` | `/api/v1/master/loan-types` | ประเภทเงินกู้ |
| `CRUD` | `/api/v1/master/loan-steps` | ขั้นตอน |
| `CRUD` | `/api/v1/master/loan-docs` | เอกสาร |
| `CRUD` | `/api/v1/master/loan-appts` | นัดหมาย |

> **Legend:** ❌ = Public | ✅ = Login | 🟡 = Officer | 🔴 = Admin

---

## 🧪 Testing Examples

### สร้าง Mortgage

```bash
curl -X POST http://localhost:3000/api/v1/mortgages \
  -H "Content-Type: application/json" \
  -b cookies.txt \
  -d '{
    "memb_no": "07337",
    "loan_type_id": 1,
    "amount": 500000,
    "purpose": "ซื้อบ้าน",
    "collateral": "โฉนดที่ดิน เลขที่ 12345"
  }'
```

### เปลี่ยนสถานะ

```bash
curl -X PUT http://localhost:3000/api/v1/mortgages/1/step \
  -H "Content-Type: application/json" \
  -b cookies.txt \
  -d '{
    "step_id": 3,
    "remark": "เอกสารครบแล้ว กำลังตรวจสอบ"
  }'
```

### อนุมัติ

```bash
curl -X PUT http://localhost:3000/api/v1/mortgages/1/approve \
  -H "Content-Type: application/json" \
  -b cookies.txt \
  -d '{
    "contract_no": "SPSC-2026-0001",
    "remark": "อนุมัติตามมติคณะกรรมการ"
  }'
```

### สร้างนัดหมาย

```bash
curl -X POST http://localhost:3000/api/v1/mortgages/1/appts \
  -H "Content-Type: application/json" \
  -b cookies.txt \
  -d '{
    "loan_appt_id": 2,
    "appt_date": "2026-01-25",
    "appt_time": "10:00",
    "remark": "นัดเซ็นสัญญา"
  }'
```

---

## 🔍 Troubleshooting

### ❌ LINE Notify ไม่ส่ง

```
ตรวจสอบ LINE_NOTIFY_TOKEN ใน .env
```

### ❌ Permission Denied

```sql
-- ตั้งเป็น OFFICER หรือ ADMIN
UPDATE users SET role = 'OFFICER' WHERE id = 1;
```

### ❌ Loan Type/Step Not Found

```
Server จะ seed ข้อมูล Master อัตโนมัติตอน start
ถ้าไม่มี ให้ restart server
```

---

## 📁 Project Structure

```
spsc-loaneasy/
├── cmd/server/main.go
├── internal/
│   ├── adapters/
│   │   ├── http/
│   │   │   ├── handlers/
│   │   │   │   ├── auth_handler.go
│   │   │   │   ├── user_handler.go
│   │   │   │   ├── mortgage_handler.go    ← Phase 4
│   │   │   │   └── master_handler.go      ← Phase 4
│   │   │   ├── middleware/
│   │   │   └── routes/routes.go
│   │   └── persistence/
│   │       ├── models/models.go           ← 10 tables
│   │       └── repositories/
│   │           ├── master_repository.go   ← Phase 4
│   │           └── mortgage_repository.go ← Phase 4
│   ├── config/
│   │   └── master_seeder.go               ← Phase 4
│   └── core/services/
│       ├── mortgage_service.go            ← Phase 4
│       └── notification_service.go        ← LINE Notify
├── .env
└── GUIDE.md
```

---

## ✅ Checklist

### Phase 4: Mortgage ✅
- [x] Master Tables (loan_types, loan_steps, loan_docs, loan_appts)
- [x] Mortgage Table (ข้อมูลหลัก)
- [x] Transaction Table (History)
- [x] Current Tables (doc, type, step, appt)
- [x] Mortgage CRUD API
- [x] Status Change API
- [x] Approve/Reject API
- [x] Document Checklist API
- [x] Appointment API
- [x] LINE Notify Service
- [x] Master Data Seeder

### Phase 5: LINE OA ⏳
- [ ] LINE User ID ใน users
- [ ] Notification Logs Table
- [ ] LINE Messaging API

---

<div align="center">

**Made with ❤️ for SPSC**

`Phase 1` ✅ → `Phase 2` ✅ → `Phase 3` ✅ → `Phase 4` ✅ → `Phase 5` ⏳

</div>
