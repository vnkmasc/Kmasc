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
