import { School, User, Book, Mail, Library, Calendar } from 'lucide-react'

export const getStudentInfoItems = (data: any) => [
  {
    icon: <School className='h-5 w-5 text-gray-500' />,
    title: 'Trường/Học viện',
    value: `${data?.universityCode} - ${data?.univeristyName}`
  },
  {
    icon: <User className='h-5 w-5 text-gray-500' />,
    title: 'Họ và tên',
    value: data?.name
  },
  {
    icon: <Book className='h-5 w-5 text-gray-500' />,
    title: 'Mã sinh viên',
    value: data?.code
  },
  {
    icon: <Mail className='h-5 w-5 text-gray-500' />,
    title: 'Email',
    value: data?.email
  },
  {
    icon: <Library className='h-5 w-5 text-gray-500' />,
    title: 'Ngành học',
    value: `${data?.facultyCode} - ${data?.facultyName}`
  },
  {
    icon: <Calendar className='h-5 w-5 text-gray-500' />,
    title: 'Năm nhập học',
    value: data?.year
  }
]
