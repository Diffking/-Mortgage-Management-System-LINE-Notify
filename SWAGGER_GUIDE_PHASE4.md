# 🔰 คู่มือทดสอบ Phase 4 ด้วย Swagger (ฉบับมือใหม่)

## 📌 Swagger คืออะไร?

Swagger คือหน้าเว็บที่แสดง API ทั้งหมด สามารถทดสอบได้เลยใน Browser ไม่ต้องติดตั้งอะไรเพิ่ม!

---

## 1️⃣ เตรียม Server

### ขั้นตอน:

```bash
# 1. เปิด Command Prompt / Terminal
cd D:\File_Claude\SPSC-loaneasy\spsc-loaneasy

# 2. ติดตั้ง swag (ครั้งแรกครั้งเดียว)
go install github.com/swaggo/swag/cmd/swag@latest

# 3. Generate Swagger docs
swag init -g cmd/server/main.go -o docs

# 4. รัน Server
set APP_MODE=dev && go run cmd/server/main.go
```

### ถ้าสำเร็จจะเห็น:
```
✅ Database connected successfully
✅ Database migration completed
✅ Master data seeded successfully
🚀 Server starting on port 3000
```

---

## 2️⃣ เปิด Swagger

### เปิด Browser แล้วไปที่:
```
http://localhost:3000/swagger/index.html
```

### จะเห็นหน้าตาแบบนี้:
```
┌────────────────────────────────────────────────────────────────┐
│  🔷 SPSC loanEasy API                              [Authorize] │
├────────────────────────────────────────────────────────────────┤
│                                                                │
│  📁 Auth                                            [▼]        │
│     POST  /api/v1/auth/register     Register user              │
│     POST  /api/v1/auth/login        Login user                 │
│     POST  /api/v1/auth/refresh      Refresh token              │
│     POST  /api/v1/auth/logout       Logout                     │
│     GET   /api/v1/auth/me           Get current user           │
│                                                                │
│  📁 Mortgages                                       [▼]        │
│     POST  /api/v1/mortgages         Create mortgage            │
│     GET   /api/v1/mortgages         List mortgages             │
│     ...                                                        │
│                                                                │
│  📁 Master                                          [▼]        │
│     ...                                                        │
│                                                                │
└────────────────────────────────────────────────────────────────┘
```

---

## 3️⃣ ทดสอบ Register (สมัครสมาชิก)

### ขั้นตอน:

1. **คลิกที่** `POST /api/v1/auth/register`

2. **คลิกปุ่ม** `Try it out` (มุมขวา)

3. **แก้ไข Request Body:**
   ```json
   {
     "memb_no": "07337",
     "username": "officer1",
     "email": "officer1@spsc.or.th",
     "password": "Password123"
   }
   ```
   
   ⚠️ **สำคัญ:** `memb_no` ต้องมีอยู่จริงใน flommast

4. **คลิกปุ่ม** `Execute` (สีน้ำเงิน)

5. **ดูผลลัพธ์:**
   ```
   ┌─────────────────────────────────────────┐
   │  Response body                          │
   ├─────────────────────────────────────────┤
   │  {                                      │
   │    "success": true,                     │
   │    "message": "Registration successful" │
   │  }                                      │
   └─────────────────────────────────────────┘
   ```

---

## 4️⃣ ทดสอบ Login (เข้าสู่ระบบ)

### ขั้นตอน:

1. **คลิกที่** `POST /api/v1/auth/login`

2. **คลิก** `Try it out`

3. **ใส่ข้อมูล:**
   ```json
   {
     "username": "officer1",
     "password": "Password123"
   }
   ```

4. **คลิก** `Execute`

5. **ดูผลลัพธ์ - จะได้ Token:**
   ```json
   {
     "success": true,
     "message": "Login successful",
     "data": {
       "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJtZW1iX25vIjoiMDczMzciLCJ1c2VybmFtZSI6Im9mZmljZXIxIiwicm9sZSI6Ik9GRklDRVIiLCJleHAiOjE3MzY1MDQ4MDB9.xxxxx",
       "user": {
         "id": 1,
         "username": "officer1",
         "role": "USER"
       }
     }
   }
   ```

6. **⭐ Copy Token:**
   - หาบรรทัด `"access_token": "eyJhbG..."`
   - Copy ค่าที่อยู่ใน `"..."` (ไม่รวมเครื่องหมาย ")

---

## 5️⃣ ใส่ Token ใน Swagger (Authorize)

### ขั้นตอน:

1. **เลื่อนขึ้นไปด้านบนสุด**

2. **คลิกปุ่ม** `Authorize` 🔓 (มุมขวาบน)

3. **จะเห็นหน้าต่าง:**
   ```
   ┌────────────────────────────────────────────────┐
   │  Available authorizations                      │
   ├────────────────────────────────────────────────┤
   │                                                │
   │  BearerAuth (http, Bearer)                     │
   │  ┌──────────────────────────────────────────┐  │
   │  │  Value: ___________________________      │  │
   │  └──────────────────────────────────────────┘  │
   │                                                │
   │  [Authorize]  [Close]                          │
   └────────────────────────────────────────────────┘
   ```

4. **พิมพ์ในช่อง Value:**
   ```
   Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjox...
   ```
   
   ⚠️ **สำคัญ:** ต้องพิมพ์ `Bearer ` (มีเว้นวรรค) ตามด้วย Token

5. **คลิก** `Authorize`

6. **คลิก** `Close`

7. **จะเห็นไอคอน 🔒 แสดงว่า Authorize แล้ว**

---

## 6️⃣ ตั้งค่า Role ใน Database

### ⚠️ สำคัญ! ก่อนทดสอบ Mortgage ต้องตั้ง Role ก่อน

1. **เปิด MySQL (phpMyAdmin, HeidiSQL, หรือ Command Line)**

2. **รัน SQL:**
   ```sql
   -- ดู users ก่อน
   SELECT id, username, role FROM users;
   
   -- ตั้งเป็น OFFICER (ทดสอบ Mortgage)
   UPDATE users SET role = 'OFFICER' WHERE username = 'officer1';
   
   -- หรือตั้งเป็น ADMIN (ทดสอบทุกอย่าง)
   UPDATE users SET role = 'ADMIN' WHERE username = 'officer1';
   ```

3. **⭐ Login ใหม่ใน Swagger** (เพื่อให้ Token มี Role ใหม่)

4. **ใส่ Token ใหม่ใน Authorize**

---

## 7️⃣ ทดสอบดู Master Data

### ทดสอบ List Loan Types:

1. **คลิก** `GET /api/v1/master/loan-types`

2. **คลิก** `Try it out`

3. **คลิก** `Execute`

4. **ผลลัพธ์:**
   ```json
   {
     "success": true,
     "data": {
       "loan_types": [
         {
           "id": 1,
           "code": "NORMAL",
           "name": "สินเชื่อสามัญ",
           "interest_rate": 6.5
         },
         {
           "id": 2,
           "code": "EMERGENCY", 
           "name": "สินเชื่อฉุกเฉิน",
           "interest_rate": 6
         },
         {
           "id": 3,
           "code": "SPECIAL",
           "name": "สินเชื่อพิเศษ",
           "interest_rate": 5.5
         }
       ]
     }
   }
   ```

### ทดสอบ List Loan Steps:

1. **คลิก** `GET /api/v1/master/loan-steps`
2. **คลิก** `Try it out` → `Execute`
3. **จะเห็น 7 ขั้นตอน:** DRAFT, PENDING_DOC, REVIEWING, PENDING_APPROVE, APPROVED, REJECTED, CANCELLED

---

## 8️⃣ ทดสอบสร้าง Mortgage

### ขั้นตอน:

1. **คลิก** `POST /api/v1/mortgages`

2. **คลิก** `Try it out`

3. **แก้ไข Body:**
   ```json
   {
     "memb_no": "07337",
     "loan_type_id": 1,
     "amount": 500000,
     "purpose": "ซื้อบ้าน",
     "collateral": "โฉนดที่ดิน เลขที่ 12345",
     "remark": "คำขอใหม่"
   }
   ```

4. **คลิก** `Execute`

5. **ผลลัพธ์:**
   ```json
   {
     "success": true,
     "message": "Mortgage created successfully",
     "data": {
       "mortgage": {
         "id": 1,
         "memb_no": "07337",
         "amount": 500000,
         "current_step_id": 1,
         "current_step_name": "ร่างคำขอ"
       }
     }
   }
   ```

---

## 9️⃣ ทดสอบเปลี่ยนสถานะ

### ขั้นตอน:

1. **คลิก** `PUT /api/v1/mortgages/{id}/step`

2. **คลิก** `Try it out`

3. **ใส่ id:** `1`

4. **แก้ไข Body:**
   ```json
   {
     "step_id": 2,
     "remark": "รอเอกสารเพิ่มเติม"
   }
   ```

5. **คลิก** `Execute`

6. **ทำซ้ำเพื่อเปลี่ยนสถานะต่อ:**
   - step_id: 3 → "กำลังตรวจสอบ"
   - step_id: 4 → "รออนุมัติ"

---

## 🔟 ทดสอบอนุมัติ

### ขั้นตอน:

1. **คลิก** `PUT /api/v1/mortgages/{id}/approve`

2. **คลิก** `Try it out`

3. **ใส่ id:** `1`

4. **แก้ไข Body:**
   ```json
   {
     "contract_no": "SPSC-2026-0001",
     "remark": "อนุมัติตามมติคณะกรรมการ"
   }
   ```

5. **คลิก** `Execute`

6. **ผลลัพธ์:**
   ```json
   {
     "success": true,
     "message": "Mortgage approved successfully"
   }
   ```

---

## 📋 ลำดับทดสอบทั้งหมด

```
┌────┬─────────────────────────────┬──────────────────────────┐
│ #  │ API                         │ หมายเหตุ                 │
├────┼─────────────────────────────┼──────────────────────────┤
│ 1  │ POST /auth/register         │ สมัครสมาชิก              │
│ 2  │ POST /auth/login            │ เข้าสู่ระบบ + Copy Token │
│ 3  │ คลิก Authorize              │ ใส่ Bearer + Token       │
│ 4  │ ตั้ง role=OFFICER ใน DB     │ SQL: UPDATE users...     │
│ 5  │ POST /auth/login (อีกครั้ง) │ ได้ Token ใหม่ที่มี Role │
│ 6  │ Authorize ใหม่              │ ใส่ Token ใหม่           │
├────┼─────────────────────────────┼──────────────────────────┤
│ 7  │ GET /master/loan-types      │ ดูประเภทเงินกู้          │
│ 8  │ GET /master/loan-steps      │ ดูขั้นตอน                │
│ 9  │ GET /master/loan-docs       │ ดูเอกสาร                 │
│ 10 │ GET /master/loan-appts      │ ดูประเภทนัดหมาย          │
├────┼─────────────────────────────┼──────────────────────────┤
│ 11 │ POST /mortgages             │ สร้างคำขอ                │
│ 12 │ GET /mortgages              │ ดูรายการ                 │
│ 13 │ GET /mortgages/1            │ ดูรายละเอียด             │
│ 14 │ GET /mortgages/1/docs       │ ดูเอกสาร                 │
│ 15 │ PUT /mortgages/1/docs       │ อัพเดทเอกสาร             │
│ 16 │ PUT /mortgages/1/step       │ เปลี่ยนสถานะ (2,3,4)     │
│ 17 │ POST /mortgages/1/appts     │ สร้างนัดหมาย             │
│ 18 │ PUT /mortgages/1/approve    │ อนุมัติ                  │
│ 19 │ GET /mortgages/1/history    │ ดูประวัติ                │
└────┴─────────────────────────────┴──────────────────────────┘
```

---

## ❌ แก้ปัญหาที่พบบ่อย

### Error: 401 Unauthorized
```
สาเหตุ: ยังไม่ได้ Authorize หรือ Token หมดอายุ
แก้ไข: 
  1. Login ใหม่
  2. Copy Token
  3. คลิก Authorize → ใส่ "Bearer xxxxx"
```

### Error: 403 Forbidden
```
สาเหตุ: Role ไม่มีสิทธิ์เข้าถึง API นี้
แก้ไข:
  1. ตั้ง role ใน Database
  2. Login ใหม่เพื่อให้ Token มี Role ใหม่
  3. Authorize ด้วย Token ใหม่
```

### Error: 404 Member Not Found
```
สาเหตุ: memb_no ไม่มีอยู่ใน flommast
แก้ไข: ใช้ memb_no ที่มีอยู่จริง เช่น 07337
```

### Error: 404 Loan Type Not Found
```
สาเหตุ: loan_type_id ไม่มี
แก้ไข: 
  1. ดู loan_type_id จาก GET /master/loan-types
  2. ใช้ id ที่มีอยู่ (1, 2, หรือ 3)
```

### Swagger ไม่ขึ้น
```
สาเหตุ: ยังไม่ได้ Generate docs
แก้ไข:
  swag init -g cmd/server/main.go -o docs
  go run cmd/server/main.go
```

---

## 💡 Tips

1. **Token หมดอายุทุก 15 นาที** → Login ใหม่เป็นระยะ

2. **ดู Response Code:**
   - 200 = สำเร็จ
   - 201 = สร้างสำเร็จ
   - 400 = ข้อมูลผิด
   - 401 = ไม่ได้ Login
   - 403 = ไม่มีสิทธิ์
   - 404 = ไม่พบข้อมูล

3. **Copy Response:** คลิกขวาที่ Response → Copy

4. **ดูเฉพาะ Active:** เพิ่ม `?all=true` เพื่อดูรวม inactive

---

## ✅ Checklist ทดสอบ Phase 4

- [ ] Server รันได้
- [ ] Swagger เปิดได้ที่ localhost:3000/swagger
- [ ] Register สำเร็จ
- [ ] Login สำเร็จ + ได้ Token
- [ ] Authorize ใส่ Token แล้ว
- [ ] ตั้ง role = OFFICER/ADMIN ใน DB
- [ ] Login ใหม่ + Authorize ใหม่
- [ ] GET /master/loan-types เห็นข้อมูล
- [ ] POST /mortgages สร้างคำขอได้
- [ ] PUT /mortgages/1/step เปลี่ยนสถานะได้
- [ ] PUT /mortgages/1/approve อนุมัติได้
- [ ] GET /mortgages/1/history เห็นประวัติ

---

🎉 **เสร็จแล้ว! คุณทดสอบ Phase 4 ได้ครบแล้ว**
