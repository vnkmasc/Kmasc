import { format } from 'date-fns'

export const formatStudent = (data: any, isSendToServer: boolean = false) => {
  return isSendToServer
    ? {
        student_code: data.code,
        full_name: data.name,
        email: data.email,
        faculty_code: data.faculty,
        course: String(data.year ?? ''),
        citizen_id_number: data.citizenId,
        ethnicity: data.ethnicity,
        current_address: data.currentAddress,
        birth_address: data.birthAddress,
        union_join_date: data.unionJoinDate,
        party_join_date: data.partyJoinDate,
        description: data.description,
        date_of_birth: data.dateOfBirth ? format(new Date(data.dateOfBirth), 'dd/MM/yyyy') : undefined,
        gender: Boolean(data.gender)
      }
    : {
        id: data.id,
        code: data.student_code,
        name: data.full_name,
        email: data.email,
        faculty: data.faculty_code,
        facultyName: data.faculty_name,
        year: data.course,
        status: String(data.status),
        citizenId: data.citizen_id_number,
        ethnicity: data.ethnicity,
        currentAddress: data.current_address,
        birthAddress: data.birth_address,
        unionJoinDate: data.union_join_date,
        partyJoinDate: data.party_join_date,
        description: data.description,
        dateOfBirth: data.date_of_birth,
        gender: String(data.gender)
      }
}

export const formatFaculty = (data: any, isSendToServer: boolean = false) => {
  return isSendToServer
    ? {
        faculty_code: data.code,
        faculty_name: data.name,
        description: data.description
      }
    : {
        id: data.id,
        code: data.faculty_code,
        name: data.faculty_name,
        description: data.description
      }
}

export const formatFacultyOptions = (data: any) => {
  return data.map((item: any) => ({
    label: item.name,
    value: item.code
  }))
}

export const formatCertificate = (data: any, isSendToServer: boolean = false) => {
  return isSendToServer
    ? {
        student_code: data.studentCode,
        name: data.name,
        certificate_type: data.certificateType ? Number(data.certificateType) : undefined,
        serial_number: data.serialNumber,
        reg_no: data.regNo,
        issue_date: data.date ? format(new Date(data.date), 'dd/MM/yyyy') : undefined
      }
    : {
        id: data.id,
        studentCode: data.student_code,
        studentName: data.student_name,
        faculty: data.faculty_code,
        facultyName: data.faculty_name,
        certificateType: data.certificate_type,
        date: data.issue_date,
        signed: data.signed,
        name: data.name,
        isDegree: data.certificate_type !== undefined
      }
}

export const formatCertificateView = (data: any) => {
  return {
    studentCode: data.student_code,
    studentName: data.student_name,
    facultyCode: data.faculty_code,
    facultyName: data.faculty_name,
    certificateType: data.certificate_type,
    date: data.issue_date,
    name: data.name,
    universityName: data.university_name,
    universityCode: data.university_code,
    serialNumber: data.serial_number,
    regNo: data.reg_no,
    signed: data.signed,
    description: data.description,
    gpa: data.gpa,
    dateOfBirth: data.date_of_birth || '',
    course: data.course || '',
    graduationRank: data.graduation_rank || '',
    major: data.major || data.faculty_name || ''
  }
}

export const formatCertificateVerifyCode = (data: any, isSendToServer: boolean = false) => {
  return isSendToServer
    ? {
        duration_minutes: data.expiredAfter,
        can_view_score: data.permissionType.includes('can_view_score'),
        can_view_data: data.permissionType.includes('can_view_data'),
        can_view_file: data.permissionType.includes('can_view_file')
      }
    : {
        verifyCode: data.code,
        createdAt: format(new Date(data.created_at), 'dd/MM/yyyy HH:mm:ss'),
        expiredAfter: data.expired_in_minutes,
        permissionType: [
          data.can_view_score ? 'can_view_score' : null,
          data.can_view_data ? 'can_view_data' : null,
          data.can_view_file ? 'can_view_file' : null
        ].filter(Boolean) as ('can_view_score' | 'can_view_data' | 'can_view_file')[],
        status: data.expired_in_minutes !== 0
      }
}

const getIsDiscipline = (data: any) => {
  switch (data) {
    case 'true':
      return true
    case 'false':
      return false
    case undefined:
      return undefined
    case '':
      return undefined
    default:
      return undefined
  }
}

export const formatRewardDiscipline = (data: any, isSendToServer: boolean = false) => {
  return isSendToServer
    ? {
        student_code: data.studentCode,
        name: data.name,
        decision_number: data.decisionNumber,
        is_discipline: !!data.disciplineLevel || getIsDiscipline(data.isDiscipline),
        description: data.description,
        discipline_level: data.disciplineLevel ? Number(data.disciplineLevel) : undefined
      }
    : {
        id: data.id,
        name: data.name,
        studentCode: data.student_code,
        studentName: data.student_name,
        faculty: data.faculty_code,
        facultyName: data.faculty_name,
        decisionNumber: data.decision_number,
        description: data.description,
        isDiscipline: data.is_discipline,
        disciplineLevel: String(data.discipline_level),
        createdAt: format(new Date(data.created_at), 'dd/MM/yyyy HH:mm:ss')
      }
}

export const getCertificatePreviewProps = (data: any) => {
  const cert = data?.certificate
  if (!cert) return null

  // Extract graduation year from course or issue date
  const graduationYear = cert.course
    ? `20${cert.course.substring(2)}`
    : cert.date
      ? new Date(cert.date).getFullYear().toString()
      : '2024'

  // Format issue date to "Hà Nội, ngày..., tháng..., năm...."
  const formatIssueDate = (dateStr: string) => {
    if (!dateStr) {
      const today = new Date()
      return `Hà Nội, ngày ${today.getDate()}, tháng ${today.getMonth() + 1}, năm ${today.getFullYear()}`
    }

    const date = new Date(dateStr)
    return `Hà Nội, ngày ${date.getDate()}, tháng ${date.getMonth() + 1}, năm ${date.getFullYear()}`
  }

  return {
    degreeType: cert.certificateType || 'Văn bằng đại học',
    major: cert.major || cert.facultyName || 'Không xác định',
    recipientName: cert.studentName || 'Không xác định',
    dateOfBirth: cert.dateOfBirth || '01/01/1990',
    graduationYear: graduationYear || '2024',
    grade: cert.graduationRank || (cert.gpa ? `GPA: ${cert.gpa}` : 'Khá'),
    issueDate: formatIssueDate(cert.date)
  }
}
