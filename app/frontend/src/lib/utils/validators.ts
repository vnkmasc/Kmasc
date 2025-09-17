import { z } from 'zod'

export const validateEmail = z
  .string()
  .trim()
  .nonempty({
    message: 'Email không được để trống'
  })
  .email({
    message: 'Email không hợp lệ (VD: example@gmail.com)'
  })

export const validatePassword = z.string().trim().nonempty({
  message: 'Mật khẩu không được để trống'
})
// .min(8, {
//   message: 'Mật khẩu phải có ít nhất 8 ký tự'
// })
export const validateAcademicEmail = z
  .string()
  .trim()
  .nonempty({
    message: 'Email không được để trống'
  })
  .email({
    message: 'Email không hợp lệ (VD: example@gmail.com)'
  })
  .includes('edu.vn', {
    message: 'Email học viện không hợp lệ (VD: example@actvn.edu.vn)'
  })

export const validateNoEmpty = (name: string) => {
  return z
    .string()
    .trim()
    .nonempty({
      message: `${name} không được để trống`
    })
}

export const validateCitizenId = z
  .string()
  .trim()
  .nonempty({
    message: 'CMND không được để trống'
  })
  .min(12, {
    message: 'CMND phải có 12 ký tự'
  })
  .regex(/^\d+$/, {
    message: 'CMND chỉ được chứa các ký tự số'
  })

export const validateGPA = z
  .number()
  .or(z.string().transform((val) => parseFloat(val)))
  .refine((val) => val >= 0, {
    message: 'Điểm GPA phải lớn hơn 0'
  })
  .refine((val) => val <= 10, {
    message: 'Điểm GPA phải nhỏ hơn 10'
  })
