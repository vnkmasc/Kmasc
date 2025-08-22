import { queryString } from '../utils/common'
import apiService from './root'

export const createDegreeTemplate = async (data: any) => {
  const res = await apiService('POST', 'templates', data)
  return res
}

export const getDegreeTemplateById = async (id: string) => {
  const res = await apiService('GET', `templates/${id}`)

  return {
    ...res.data,
    faculty_id: res.data.facultyId
  }
}

export const updateDegreeTemplate = async (id: string, data: any) => {
  const res = await apiService('PUT', `templates/${id}`, data)
  return res
}
export const searchDegreeTemplateByFaculty = async (facultyId: string) => {
  const res = await apiService('GET', `templates/faculty?faculty_id=${facultyId}`)
  return res
}

export const deleteDegreeTemplate = async (id: string) => {
  const res = await apiService('DELETE', `templates/${id}`)
  return res
}

export const signDegreeTemplateFaculty = async (facultyId: string) => {
  const res = await apiService('POST', `templates/sign/faculty/${facultyId}`)
  return res
}

export const signDegreeTemplateUni = async () => {
  const res = await apiService('POST', 'templates/sign/university')
  return res
}

export const signDegreeTemplateById = async (id: string, signature: string) => {
  const res = await apiService('POST', `templates/${id}/sign`, { signature })
  return res
}

export const searchDigitalDegreeList = async (params: any) => {
  const res = await apiService('GET', queryString(['ediplomas', 'search'], params))
  return res
}

export const issueDownloadDegreeZip = async (facultyId: string, templateId: string) => {
  const blob = await apiService(
    'POST',
    `ediplomas/generate-bulk-zip`,
    {
      faculty_id: facultyId,
      template_id: templateId
    },
    true,
    { Accept: 'application/zip' },
    true
  )

  return blob
}

export const uploadDigitalDegreesMinio = async (data: FormData) => {
  const res = await apiService('POST', 'ediplomas/upload-zip', data)
  return res
}

export const uploadDigitalDegreesBlockchain = async (data: FormData) => {
  const res = await apiService('POST', 'blockchain/push-ediploma', data)
  return res
}

export const getTemplateInterfaces = async () => {
  const res = await apiService('GET', 'template-samples')
  return res
}

export const getTemplateInterfaceById = async (id: string) => {
  const res = await apiService('GET', `template-samples/${id}`)
  return res
}

export const updateTemplateInterface = async (id: string, data: any) => {
  const res = await apiService('PUT', `template-samples/${id}`, data)
  return res
}

export const createTemplateInterface = async (data: any) => {
  const res = await apiService('POST', 'template-samples', data)
  return res
}

export const getDigitalDegreeFileById = async (id: string) => {
  const res = await apiService('GET', `ediplomas/file/${id}`, undefined, true, {}, true)
  return res
}

export const getDigitalDegreeById = async (id: string) => {
  const res = await apiService('GET', `ediplomas/${id}`)
  return res
}

export const verifyDigitalDegreeDataBlockchain = async (
  universityId: string,
  facultyId: string,
  certificateType: string,
  course: string,
  ediplomaId: string
) => {
  const formData = new FormData()
  formData.append('university_id', universityId)
  if (facultyId !== '') formData.append('faculty_id', facultyId)
  if (certificateType !== '') formData.append('certificate_type', certificateType)
  if (course !== '') formData.append('course', course)
  if (ediplomaId !== '') formData.append('ediploma_id', ediplomaId)

  // const res = await apiService('POST', 'blockchain/verify-batch', formData, false)
  return {
    data: {
      id: '68a77870aacf9ece849c85f0',
      certificate_id: '68a77870aacf9ece849c85ef',
      name: 'Bằng Cử nhân Khoa học máy tính',
      template_name: 'Bằng mẫu',
      university_code: 'KMA',
      university_name: 'Học viện Kỹ thuật Mật mã',
      faculty_id: '68a1cdf3688b235903145551',
      faculty_code: 'CNTT',
      faculty_name: 'Công nghệ thông tin',
      student_name: 'Ho Ngoc Yen',
      student_code: 'SV000011',
      full_name: 'Ho Ngoc Yen',
      certificate_type: 'Kỹ sư',
      course: '2025',
      education_type: 'Chính quy',
      gpa: 3.94,
      graduation_rank: 'Xuất sắc',
      issue_date: '01/01/0001',
      serial_number: '10093',
      registration_number: '10093',
      issued: true,
      signed: false,
      data_encrypted: true,
      on_blockchain: true
    },
    message: 'Dữ liệu khớp hoàn toàn trên chuỗi khối',
    verified: true
  }
}

export const verifyDigitalDegreeFileBlockchain = async (universityCode: string, ediplomaId: string) => {
  const res = await apiService(
    'GET',
    queryString(['auth', 'ediploma', 'file'], { university_code: universityCode, ediploma_id: ediplomaId }),
    undefined,
    true,
    {},
    true
  )

  return res
}
