# 📘 SPSC loanEasy v1.0 - Phase 4 Design Guide

> ระบบจำนอง (Mortgage System) + LINE Notification

```
╔═══════════════════════════════════════════════════════════════╗
║                    🏦 Phase 4: Mortgage System                ║
║                                                               ║
║  📊 10 Tables    📡 25+ APIs    🔔 LINE Notify               ║
╚═══════════════════════════════════════════════════════════════╝
```

---

## 📋 สารบัญ

| หัวข้อ | คำอธิบาย |
|--------|----------|
| [🎯 Overview](#-overview) | ภาพรวมระบบ |
| [📊 ERD](#-erd) | Entity Relationship Diagram |
| [🔵 Master Tables](#-master-tables) | ตาราง Reference (4 ตาราง) |
| [🟢 Main Tables](#-main-tables) | ตารางหลัก (2 ตาราง) |
| [🟡 Current Tables](#-current-tables) | ตารางสถานะปัจจุบัน (4 ตาราง) |
| [📡 API Endpoints](#-api-endpoints) | รายการ API ทั้งหมด |
| [🔄 Flow](#-flow) | การทำงานของระบบ |
| [🔔 LINE Notification](#-line-notification) | ระบบแจ้งเตือน |
| [📁 Project Structure](#-project-structure) | โครงสร้างไฟล์ |

---

## 🎯 Overview

### จุดประสงค์หลัก

```
┌─────────────────────────────────────────────────────────────────┐
│  🎯 เป้าหมาย Phase 4                                            │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  1. ระบบจำนอง (Mortgage)                                        │
│     └── สมาชิกดูข้อมูลได้ (READ ONLY)                           │
│     └── เจ้าหน้าที่ Create/Update                               │
│                                                                 │
│  2. เก็บ History ทุกการเปลี่ยนแปลง                               │
│     └── ตาราง transaction                                       │
│     └── ทำลายทุก 1 ปี                                           │
│                                                                 │
│  3. LINE แจ้งเตือนอัตโนมัติ                                      │
│     └── แจ้งกลุ่มเจ้าหน้าที่                                     │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### Roles & Permissions

| Role | mortgage | Current Tables | Master Tables | History |
|------|----------|----------------|---------------|---------|
| 🟢 USER | ดูของตัวเอง | ❌ | ❌ | ดูของตัวเอง |
| 🟡 OFFICER | CRUD ทั้งหมด | Update | ❌ | ดูทั้งหมด |
| 🔴 ADMIN | CRUD ทั้งหมด | Update | CRUD | ดูทั้งหมด |

---

## 📊 ERD

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              📊 ERD - Phase 4                               │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   ┌─────────────────────────────────────────────────────────────────────┐  │
│   │                           mortgage                                   │  │
│   │                    (ตารางหลัก - สมาชิกดูได้)                          │  │
│   │                                                                      │  │
│   │  FK: loan_type_id ────────────────────────────► loan_type (Master)  │  │
│   │  FK: current_step_id ─────────────────────────► loan_step (Master)  │  │
│   │  FK: current_appt_id ─────────────────────────► loan_appt_current   │  │
│   │  FK: officer_id ──────────────────────────────► users               │  │
│   │  FK: memb_no ─────────────────────────────────► flommast            │  │
│   └──────────────────────────┬──────────────────────────────────────────┘  │
│                              │                                             │
│              ┌───────────────┴───────────────┐                             │
│              │ 1:N                     1:N   │                             │
│              ▼                               ▼                             │
│   ┌─────────────────────┐         ┌─────────────────────┐                 │
│   │  loan_doc_current   │         │     transaction     │                 │
│   │   (เอกสาร 1:N)      │         │   (History 1:N)     │                 │
│   │                     │         │                     │                 │
│   │  FK: loan_doc_id ───┼────────►│  loan_doc (Master)  │                 │
│   └─────────────────────┘         └──────────┬──────────┘                 │
│                                              │                             │
│                          ┌───────────────────┼───────────────────┐        │
│                          │                   │                   │        │
│                          ▼                   ▼                   ▼        │
│                 ┌──────────────┐    ┌──────────────┐    ┌──────────────┐ │
│                 │loan_type_    │    │loan_step_    │    │loan_appt_    │ │
│                 │  current     │    │  current     │    │  current     │ │
│                 └──────┬───────┘    └──────┬───────┘    └──────┬───────┘ │
│                        │                   │                   │         │
│                        ▼                   ▼                   ▼         │
│   ┌─────────────────────────────────────────────────────────────────────┐│
│   │                         Master Tables                                ││
│   ├────────────────┬────────────────┬────────────────┬──────────────────┤│
│   │   loan_type    │   loan_step    │   loan_doc     │    loan_appt     ││
│   │  (ประเภทกู้)   │   (ขั้นตอน)    │   (เอกสาร)     │   (ประเภทนัด)    ││
│   └────────────────┴────────────────┴────────────────┴──────────────────┘│
│                                                                           │
└───────────────────────────────────────────────────────────────────────────┘
```

---

## 🔵 Master Tables

> ข้อมูลคงที่ใช้ Reference - Admin จัดการ

### 1️⃣ loan_type (ประเภทเงินกู้)

| Field | Type | Null | Description |
|-------|------|:----:|-------------|
| id | UINT | NO | รหัส (PK, Auto) |
| code | VARCHAR(20) | NO | รหัสประเภท (unique) |
| name | VARCHAR(100) | NO | ชื่อประเภท |
| description | TEXT | YES | รายละเอียด |
| interest_rate | DECIMAL(5,2) | NO | อัตราดอกเบี้ย % ต่อปี |
| is_active | BOOL | NO | ใช้งานอยู่ (default: true) |
| created_at | DATETIME | NO | วันที่สร้าง |
| updated_at | DATETIME | NO | วันที่แก้ไข |

**ตัวอย่างข้อมูล:**

| id | code | name | interest_rate |
|:--:|------|------|:-------------:|
| 1 | NORMAL | สินเชื่อสามัญ | 6.50 |
| 2 | EMERGENCY | สินเชื่อฉุกเฉิน | 6.00 |
| 3 | SPECIAL | สินเชื่อพิเศษ | 5.50 |

---

### 2️⃣ loan_step (ขั้นตอน/สถานะ)

| Field | Type | Null | Description |
|-------|------|:----:|-------------|
| id | UINT | NO | รหัส (PK, Auto) |
| code | VARCHAR(20) | NO | รหัสขั้นตอน (unique) |
| name | VARCHAR(100) | NO | ชื่อขั้นตอน |
| description | TEXT | YES | รายละเอียด |
| step_order | INT | NO | ลำดับขั้นตอน |
| color | VARCHAR(20) | YES | สีแสดงผล (#hex) |
| is_final | BOOL | NO | เป็นขั้นตอนสุดท้าย |
| is_active | BOOL | NO | ใช้งานอยู่ |
| created_at | DATETIME | NO | วันที่สร้าง |
| updated_at | DATETIME | NO | วันที่แก้ไข |

**ตัวอย่างข้อมูล:**

| id | code | name | step_order | color | is_final |
|:--:|------|------|:----------:|-------|:--------:|
| 1 | DRAFT | ร่างคำขอ | 1 | #6B7280 | false |
| 2 | PENDING_DOC | รอเอกสาร | 2 | #F59E0B | false |
| 3 | REVIEWING | กำลังตรวจสอบ | 3 | #3B82F6 | false |
| 4 | PENDING_APPROVE | รออนุมัติ | 4 | #8B5CF6 | false |
| 5 | APPROVED | อนุมัติแล้ว | 5 | #10B981 | true |
| 6 | REJECTED | ปฏิเสธ | 6 | #EF4444 | true |
| 7 | CANCELLED | ยกเลิก | 7 | #6B7280 | true |

---

### 3️⃣ loan_doc (ประเภทเอกสาร)

> ⚠️ แค่บอกว่าต้องเตรียมอะไร - ไม่มีการ Upload ไฟล์

| Field | Type | Null | Description |
|-------|------|:----:|-------------|
| id | UINT | NO | รหัส (PK, Auto) |
| code | VARCHAR(20) | NO | รหัสเอกสาร (unique) |
| name | VARCHAR(100) | NO | ชื่อเอกสาร |
| description | TEXT | YES | รายละเอียด |
| is_active | BOOL | NO | ใช้งานอยู่ (default: true) |
| created_at | DATETIME | NO | วันที่สร้าง |
| updated_at | DATETIME | NO | วันที่แก้ไข |

**ตัวอย่างข้อมูล:**

| id | code | name |
|:--:|------|------|
| 1 | ID_CARD | สำเนาบัตรประชาชน |
| 2 | HOUSE_REG | สำเนาทะเบียนบ้าน |
| 3 | SALARY_SLIP | สลิปเงินเดือน 3 เดือนล่าสุด |
| 4 | LAND_TITLE | โฉนดที่ดิน (กรณีค้ำประกัน) |
| 5 | GUARANTOR_ID | บัตรผู้ค้ำประกัน |

---

### 4️⃣ loan_appt (ประเภทนัดหมาย)

> ⚠️ เก็บแค่ "ประเภท" - วันที่จริงอยู่ใน loan_appt_current

| Field | Type | Null | Description |
|-------|------|:----:|-------------|
| id | UINT | NO | รหัส (PK, Auto) |
| code | VARCHAR(20) | NO | รหัสนัดหมาย (unique) |
| name | VARCHAR(100) | NO | ชื่อประเภทนัดหมาย |
| description | TEXT | YES | รายละเอียด |
| default_location | VARCHAR(200) | YES | สถานที่เริ่มต้น |
| is_active | BOOL | NO | ใช้งานอยู่ |
| created_at | DATETIME | NO | วันที่สร้าง |
| updated_at | DATETIME | NO | วันที่แก้ไข |

**ตัวอย่างข้อมูล:**

| id | code | name | default_location |
|:--:|------|------|------------------|
| 1 | SIGN_CONTRACT | นัดเซ็นสัญญา | ห้องประชุม สหกรณ์ฯ |
| 2 | SUBMIT_DOC | นัดส่งเอกสาร | เคาน์เตอร์บริการ |
| 3 | RECEIVE_MONEY | นัดรับเงิน | ห้องการเงิน |

---

## 🟢 Main Tables

### 5️⃣ mortgage (ข้อมูลจำนอง) ⭐

> ตารางหลัก - สมาชิกดูได้ (READ ONLY) - เจ้าหน้าที่ Create/Update

| Field | Type | Null | Description |
|-------|------|:----:|-------------|
| id | UINT | NO | รหัส (PK, Auto) |
| contract_no | VARCHAR(50) | YES | เลขที่สัญญา (unique, หลังอนุมัติ) |
| memb_no | VARCHAR(20) | NO | เลขสมาชิก (FK → flommast) |
| officer_id | UINT | NO | เจ้าหน้าที่ผู้รับผิดชอบ (FK → users) |
| user_id | UINT | NO | ผู้สร้าง (FK → users) |
| amount | DECIMAL(15,2) | NO | จำนวนเงินกู้ |
| collateral | TEXT | YES | หลักประกัน/ทรัพย์สินค้ำประกัน |
| purpose | TEXT | YES | วัตถุประสงค์การกู้ |
| guarantor_memb_no | VARCHAR(20) | YES | เลขสมาชิกผู้ค้ำประกัน |
| loan_type_id | UINT | NO | ประเภทเงินกู้ปัจจุบัน (FK → loan_type) |
| interest_rate | DECIMAL(5,2) | NO | ดอกเบี้ยปัจจุบัน % |
| current_step_id | UINT | NO | ขั้นตอนปัจจุบัน (FK → loan_step) |
| current_appt_id | UINT | YES | นัดหมายปัจจุบัน (FK → loan_appt_current) |
| approved_by | UINT | YES | ผู้อนุมัติ (FK → users) |
| approved_at | DATETIME | YES | วันที่อนุมัติ |
| remark | TEXT | YES | หมายเหตุ |
| created_at | DATETIME | NO | วันที่สร้าง |
| updated_at | DATETIME | NO | วันที่แก้ไข |
| deleted_at | DATETIME | YES | วันที่ลบ (soft delete) |

---

### 6️⃣ transaction (ธุรกรรม/History) ⭐

> บันทึกทุกการเปลี่ยนแปลง - ไม่มีการลบ - ทำลายทุก 1 ปี

| Field | Type | Null | Description |
|-------|------|:----:|-------------|
| id | UINT | NO | รหัส (PK, Auto) |
| mortgage_id | UINT | NO | FK → mortgage |
| transaction_type | VARCHAR(50) | NO | ประเภทธุรกรรม |
| from_step_id | UINT | YES | สถานะก่อนหน้า (FK → loan_step) |
| to_step_id | UINT | YES | สถานะใหม่ (FK → loan_step) |
| from_doc_id | UINT | YES | เอกสารก่อนหน้า (FK → loan_doc) |
| to_doc_id | UINT | YES | เอกสารใหม่ (FK → loan_doc) |
| from_type_id | UINT | YES | ประเภทก่อนหน้า (FK → loan_type) |
| to_type_id | UINT | YES | ประเภทใหม่ (FK → loan_type) |
| from_appt_id | UINT | YES | นัดหมายก่อนหน้า (FK → loan_appt) |
| to_appt_id | UINT | YES | นัดหมายใหม่ (FK → loan_appt) |
| amount | DECIMAL(15,2) | YES | จำนวนเงิน (ถ้ามี) |
| description | TEXT | YES | รายละเอียด |
| performed_by | UINT | NO | ผู้ดำเนินการ (FK → users) |
| ip_address | VARCHAR(50) | YES | IP Address |
| created_at | DATETIME | NO | วันที่สร้าง |

**Transaction Types:**

| Type | Description |
|------|-------------|
| CREATE | สร้างคำขอใหม่ |
| UPDATE | แก้ไขข้อมูล |
| STATUS_CHANGE | เปลี่ยนสถานะ |
| TYPE_CHANGE | เปลี่ยนประเภทเงินกู้ |
| DOC_CHECK | ตรวจสอบเอกสาร |
| APPT_CREATE | สร้างนัดหมาย |
| APPT_COMPLETE | นัดหมายเสร็จสิ้น |
| APPT_CANCEL | ยกเลิกนัดหมาย |
| OFFICER_CHANGE | เปลี่ยนเจ้าหน้าที่ |
| APPROVE | อนุมัติ |
| REJECT | ปฏิเสธ |

---

## 🟡 Current Tables

> เจ้าหน้าที่ Update → บันทึกไป mortgage + transaction

### 7️⃣ loan_doc_current (Checklist เอกสาร)

> 1 mortgage มีหลายเอกสารได้ (1:N)

| Field | Type | Null | Description |
|-------|------|:----:|-------------|
| id | UINT | NO | รหัส (PK, Auto) |
| mortgage_id | UINT | NO | FK → mortgage |
| loan_doc_id | UINT | NO | FK → loan_doc (Master) |
| is_submitted | BOOL | NO | ส่งเอกสารแล้ว (default: false) |
| checked_by | UINT | YES | ผู้ตรวจสอบ (FK → users) |
| checked_at | DATETIME | YES | วันที่ตรวจสอบ |
| remark | TEXT | YES | หมายเหตุ |
| created_at | DATETIME | NO | วันที่สร้าง |
| updated_at | DATETIME | NO | วันที่แก้ไข |

---

### 8️⃣ loan_type_current (ประเภทเงินกู้ ณ เวลานั้น)

| Field | Type | Null | Description |
|-------|------|:----:|-------------|
| id | UINT | NO | รหัส (PK, Auto) |
| transaction_id | UINT | NO | FK → transaction |
| mortgage_id | UINT | NO | FK → mortgage |
| loan_type_id | UINT | NO | FK → loan_type (Master) |
| interest_rate | DECIMAL(5,2) | NO | อัตราดอกเบี้ย ณ ตอนนั้น |
| remark | TEXT | YES | หมายเหตุ |
| created_at | DATETIME | NO | วันที่สร้าง |

---

### 9️⃣ loan_step_current (สถานะ ณ เวลานั้น)

| Field | Type | Null | Description |
|-------|------|:----:|-------------|
| id | UINT | NO | รหัส (PK, Auto) |
| transaction_id | UINT | NO | FK → transaction |
| mortgage_id | UINT | NO | FK → mortgage |
| loan_step_id | UINT | NO | FK → loan_step (Master) |
| from_step_id | UINT | YES | สถานะก่อนหน้า |
| changed_by | UINT | NO | ผู้เปลี่ยน (FK → users) |
| changed_at | DATETIME | NO | วันที่เปลี่ยน |
| remark | TEXT | YES | หมายเหตุ |
| created_at | DATETIME | NO | วันที่สร้าง |

---

### 🔟 loan_appt_current (นัดหมาย)

> วันที่นัดจริงอยู่ตรงนี้ - default จาก officer_id

| Field | Type | Null | Description |
|-------|------|:----:|-------------|
| id | UINT | NO | รหัส (PK, Auto) |
| transaction_id | UINT | NO | FK → transaction |
| mortgage_id | UINT | NO | FK → mortgage |
| loan_appt_id | UINT | NO | FK → loan_appt (Master ประเภท) |
| appt_date | DATE | NO | วันที่นัดหมาย (default: วันปัจจุบัน) |
| appt_time | TIME | YES | เวลานัดหมาย |
| location | VARCHAR(200) | YES | สถานที่ |
| appt_by | UINT | NO | ผู้นัดหมาย (FK → users) default: officer_id |
| status | VARCHAR(20) | NO | สถานะ (PENDING/COMPLETED/CANCELLED) |
| completed_at | DATETIME | YES | วันที่เสร็จ |
| remark | TEXT | YES | หมายเหตุ |
| created_at | DATETIME | NO | วันที่สร้าง |
| updated_at | DATETIME | NO | วันที่แก้ไข |

---

## 📡 API Endpoints

### 👤 สมาชิก (USER)

| Method | Endpoint | Description |
|:------:|----------|-------------|
| `GET` | `/api/v1/mortgages/my` | ดูจำนองของตัวเอง |
| `GET` | `/api/v1/mortgages/my/:id` | ดูรายละเอียด |
| `GET` | `/api/v1/mortgages/my/:id/history` | ดู History |
| `GET` | `/api/v1/mortgages/my/:id/docs` | ดูรายการเอกสาร |
| `GET` | `/api/v1/mortgages/my/:id/appts` | ดูนัดหมาย |

### 👨‍💼 เจ้าหน้าที่ (OFFICER)

| Method | Endpoint | Description |
|:------:|----------|-------------|
| `GET` | `/api/v1/mortgages` | ดูทั้งหมด |
| `GET` | `/api/v1/mortgages/:id` | ดูรายละเอียด |
| `POST` | `/api/v1/mortgages` | สร้าง mortgage |
| `PUT` | `/api/v1/mortgages/:id` | แก้ไขข้อมูล |
| `DELETE` | `/api/v1/mortgages/:id` | ลบ (soft delete) |
| `GET` | `/api/v1/mortgages/:id/history` | ดู History |
| | | |
| `PUT` | `/api/v1/mortgages/:id/step` | เปลี่ยนสถานะ |
| `PUT` | `/api/v1/mortgages/:id/type` | เปลี่ยนประเภท |
| `PUT` | `/api/v1/mortgages/:id/officer` | เปลี่ยนเจ้าหน้าที่ |
| `PUT` | `/api/v1/mortgages/:id/approve` | อนุมัติ |
| `PUT` | `/api/v1/mortgages/:id/reject` | ปฏิเสธ |
| | | |
| `GET` | `/api/v1/mortgages/:id/docs` | ดูเอกสาร |
| `PUT` | `/api/v1/mortgages/:id/docs/:docId` | เช็คเอกสาร |
| | | |
| `GET` | `/api/v1/mortgages/:id/appts` | ดูนัดหมาย |
| `POST` | `/api/v1/mortgages/:id/appts` | สร้างนัดหมาย |
| `PUT` | `/api/v1/mortgages/:id/appts/:apptId` | แก้ไขนัดหมาย |
| `PUT` | `/api/v1/mortgages/:id/appts/:apptId/complete` | เสร็จสิ้นนัดหมาย |
| `PUT` | `/api/v1/mortgages/:id/appts/:apptId/cancel` | ยกเลิกนัดหมาย |

### 🔴 Admin - Master Tables

| Method | Endpoint | Description |
|:------:|----------|-------------|
| `GET` | `/api/v1/master/loan-types` | รายการประเภท |
| `POST` | `/api/v1/master/loan-types` | สร้างประเภท |
| `GET` | `/api/v1/master/loan-types/:id` | ดูประเภท |
| `PUT` | `/api/v1/master/loan-types/:id` | แก้ไขประเภท |
| `DELETE` | `/api/v1/master/loan-types/:id` | ลบประเภท |
| | | |
| `GET` | `/api/v1/master/loan-steps` | รายการขั้นตอน |
| `POST` | `/api/v1/master/loan-steps` | สร้างขั้นตอน |
| `GET` | `/api/v1/master/loan-steps/:id` | ดูขั้นตอน |
| `PUT` | `/api/v1/master/loan-steps/:id` | แก้ไขขั้นตอน |
| `DELETE` | `/api/v1/master/loan-steps/:id` | ลบขั้นตอน |
| | | |
| `GET` | `/api/v1/master/loan-docs` | รายการเอกสาร |
| `POST` | `/api/v1/master/loan-docs` | สร้างเอกสาร |
| `GET` | `/api/v1/master/loan-docs/:id` | ดูเอกสาร |
| `PUT` | `/api/v1/master/loan-docs/:id` | แก้ไขเอกสาร |
| `DELETE` | `/api/v1/master/loan-docs/:id` | ลบเอกสาร |
| | | |
| `GET` | `/api/v1/master/loan-appts` | รายการนัดหมาย |
| `POST` | `/api/v1/master/loan-appts` | สร้างนัดหมาย |
| `GET` | `/api/v1/master/loan-appts/:id` | ดูนัดหมาย |
| `PUT` | `/api/v1/master/loan-appts/:id` | แก้ไขนัดหมาย |
| `DELETE` | `/api/v1/master/loan-appts/:id` | ลบนัดหมาย |

---

## 🔄 Flow

### Flow หลัก

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                            🔄 Main Flow                                     │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  1️⃣ เจ้าหน้าที่สร้าง Mortgage                                               │
│     └── POST /mortgages                                                    │
│     └── กำหนด: memb_no, amount, loan_type, officer_id                      │
│     └── ระบบ: สร้าง loan_doc_current ทุกรายการ (is_submitted=false)        │
│     └── ระบบ: บันทึก transaction (CREATE)                                   │
│     └── ระบบ: แจ้ง LINE                                                     │
│                                                                             │
│  2️⃣ เจ้าหน้าที่ตรวจเอกสาร                                                   │
│     └── PUT /mortgages/:id/docs/:docId                                     │
│     └── อัพเดท: loan_doc_current.is_submitted = true                       │
│     └── ระบบ: บันทึก transaction (DOC_CHECK)                                │
│                                                                             │
│  3️⃣ เจ้าหน้าที่เปลี่ยนสถานะ                                                 │
│     └── PUT /mortgages/:id/step                                            │
│     └── อัพเดท: mortgage.current_step_id                                   │
│     └── สร้าง: loan_step_current                                           │
│     └── ระบบ: บันทึก transaction (STATUS_CHANGE)                           │
│     └── ระบบ: แจ้ง LINE                                                     │
│                                                                             │
│  4️⃣ เจ้าหน้าที่สร้างนัดหมาย                                                 │
│     └── POST /mortgages/:id/appts                                          │
│     └── สร้าง: loan_appt_current (appt_by = officer_id default)            │
│     └── อัพเดท: mortgage.current_appt_id                                   │
│     └── ระบบ: บันทึก transaction (APPT_CREATE)                             │
│     └── ระบบ: แจ้ง LINE                                                     │
│                                                                             │
│  5️⃣ เจ้าหน้าที่อนุมัติ/ปฏิเสธ                                               │
│     └── PUT /mortgages/:id/approve หรือ /reject                            │
│     └── อัพเดท: mortgage.approved_by, approved_at                          │
│     └── อัพเดท: mortgage.current_step_id = APPROVED/REJECTED               │
│     └── ระบบ: บันทึก transaction (APPROVE/REJECT)                          │
│     └── ระบบ: แจ้ง LINE                                                     │
│                                                                             │
│  6️⃣ สมาชิกดูข้อมูล (READ ONLY)                                              │
│     └── GET /mortgages/my                                                  │
│     └── GET /mortgages/my/:id                                              │
│     └── GET /mortgages/my/:id/history                                      │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Flow สถานะ

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         📊 Status Flow                                      │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   ┌─────────┐    ┌─────────────┐    ┌───────────┐    ┌────────────────┐   │
│   │  DRAFT  │ ─► │ PENDING_DOC │ ─► │ REVIEWING │ ─► │ PENDING_APPROVE│   │
│   │ ร่างคำขอ │    │  รอเอกสาร   │    │ ตรวจสอบ   │    │   รออนุมัติ     │   │
│   └─────────┘    └─────────────┘    └───────────┘    └───────┬────────┘   │
│                                                              │            │
│                           ┌──────────────────────────────────┼───────┐    │
│                           │                                  │       │    │
│                           ▼                                  ▼       ▼    │
│                    ┌────────────┐                    ┌──────────┐         │
│                    │  APPROVED  │                    │ REJECTED │         │
│                    │ อนุมัติแล้ว │                    │  ปฏิเสธ  │         │
│                    └────────────┘                    └──────────┘         │
│                                                                           │
│   * ทุกสถานะสามารถไป CANCELLED ได้                                        │
│                                                                           │
└───────────────────────────────────────────────────────────────────────────┘
```

---

## 🔔 LINE Notification

### Configuration

```env
# .env
LINE_NOTIFY_TOKEN=your_line_notify_token
LINE_NOTIFY_ENABLED=true
```

### Events ที่แจ้งเตือน

| Event | Message |
|-------|---------|
| CREATE | 📝 มีคำขอกู้ใหม่ #{id} - {memb_name} จำนวน {amount} บาท |
| STATUS_CHANGE | 🔄 คำขอ #{id} เปลี่ยนสถานะ: {from} → {to} |
| APPROVE | ✅ คำขอ #{id} อนุมัติแล้ว โดย {approver} |
| REJECT | ❌ คำขอ #{id} ถูกปฏิเสธ โดย {approver} |
| APPT_CREATE | 📅 นัดหมาย #{id} - {appt_type} วันที่ {date} |
| APPT_REMINDER | ⏰ แจ้งเตือน: พรุ่งนี้มีนัด {appt_type} - {memb_name} |

### Service Structure

```go
// internal/pkg/notification/line_notify.go

type LineNotifyService struct {
    token   string
    enabled bool
}

func (s *LineNotifyService) Send(message string) error
func (s *LineNotifyService) NotifyNewMortgage(m *Mortgage) error
func (s *LineNotifyService) NotifyStatusChange(m *Mortgage, from, to string) error
func (s *LineNotifyService) NotifyApproval(m *Mortgage, approved bool) error
func (s *LineNotifyService) NotifyAppointment(appt *LoanApptCurrent) error
```

---

## 📁 Project Structure

```
internal/
├── adapters/
│   ├── http/
│   │   ├── handlers/
│   │   │   ├── mortgage_handler.go          ← NEW
│   │   │   ├── master_handler.go            ← NEW (loan_type, step, doc, appt)
│   │   │   └── ...
│   │   └── routes/
│   │       └── routes.go                    ← UPDATE
│   │
│   └── persistence/
│       ├── models/
│       │   └── models.go                    ← UPDATE (10 new models)
│       └── repositories/
│           ├── mortgage_repository.go       ← NEW
│           ├── transaction_repository.go    ← NEW
│           ├── loan_type_repository.go      ← NEW
│           ├── loan_step_repository.go      ← NEW
│           ├── loan_doc_repository.go       ← NEW
│           ├── loan_appt_repository.go      ← NEW
│           └── ...
│
├── core/
│   └── services/
│       ├── mortgage_service.go              ← NEW
│       ├── master_service.go                ← NEW
│       └── ...
│
└── pkg/
    └── notification/
        └── line_notify.go                   ← NEW
```

---

## ✅ Summary

```
┌─────────────────────────────────────────────────────────────────┐
│                    📊 Phase 4 Summary                           │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  📦 Tables: 10                                                  │
│     ├── Master: 4 (loan_type, loan_step, loan_doc, loan_appt)  │
│     ├── Main: 2 (mortgage, transaction)                        │
│     └── Current: 4 (doc, type, step, appt)                     │
│                                                                 │
│  📡 APIs: 30+                                                   │
│     ├── User: 5 (ดูของตัวเอง)                                   │
│     ├── Officer: 15+ (CRUD + workflow)                         │
│     └── Admin: 20 (Master CRUD)                                │
│                                                                 │
│  🔔 LINE Notify: Service only                                  │
│     └── ไม่ต้องเพิ่มตาราง                                       │
│                                                                 │
│  📁 Files: ~15 new files                                       │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

---

<div align="center">

**Phase 4 Design Guide**

`Phase 1` ✅ → `Phase 2` ✅ → `Phase 3` ✅ → `Phase 4` 📋

</div>
