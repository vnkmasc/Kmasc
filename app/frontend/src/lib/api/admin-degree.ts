import apiService from './root'

export const getUniversities = async () => {
  const res = await apiService('GET', 'universities')
  return res.data
}

export const getFacultiesByUniversity = async (uniId: string) => {
  const res = await apiService('GET', `faculties/university/${uniId}`)
  return res.data
}

export const getTemplatesByUniAndFaculty = async (universityId: string, facultyId: string) => {
  const res = await apiService('GET', `templates/university/${universityId}/faculty/${facultyId}`)
  return res.data
}

export const signTemplateByMinedu = async (templateId: string, signature: string) => {
  const res = await apiService('POST', `templates/sign/minedu/${templateId}`, { signature })
  return res
}
