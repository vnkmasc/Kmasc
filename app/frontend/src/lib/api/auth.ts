import apiService from './root'

export const sendOTP = async (email: string) => {
  const res = await apiService('POST', 'auth/request-otp', { student_email: email }, false)
  return res
}

export const verifyOTP = async (email: string, otp: string) => {
  const res = await apiService('POST', 'auth/verify-otp', { student_email: email, otp }, false)
  return res
}

export const registerAccount = async (email: string, password: string, userId: string) => {
  const res = await apiService(
    'POST',
    'auth/register',
    {
      personal_email: email,
      password,
      user_id: userId
    },
    false
  )
  return res
}

export const requestEducationSignUp = async (email: string, name: string, code: string, address: string) => {
  const res = await apiService(
    'POST',
    'universities',
    { email_domain: email, university_name: name, university_code: code, address: address },
    false
  )
  return res
}

export const changePassword = async (oldPassword: string, newPassword: string) => {
  const res = await apiService('POST', 'auth/change-password', { old_password: oldPassword, new_password: newPassword })
  return res
}
