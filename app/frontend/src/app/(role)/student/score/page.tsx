import TableList from '@/components/role/education-admin/table-list'
import ScoreView from '@/components/common/score-view'
import { Separator } from '@/components/ui/separator'
import { Badge } from '@/components/ui/badge'

const StudentScorePage = () => {
  const data = [
    {
      id: '1',
      name: 'Toán',
      score1: 8,
      score2: 9,
      examScore: 10,
      totalScore: 10,
      letterScore: 'A'
    },
    {
      id: '2',
      name: 'Vật lý',
      score1: 7,
      score2: 8,
      examScore: 8,
      totalScore: 9,
      letterScore: 'B'
    },
    {
      id: '3',
      name: 'Hóa học',
      score1: 9,
      score2: 9,
      examScore: 9,
      totalScore: 5,
      letterScore: 'A'
    },
    {
      id: '4',
      name: 'Sinh học',
      score1: 6,
      score2: 7,
      examScore: 7,
      totalScore: 8,
      letterScore: 'C'
    },
    {
      id: '5',
      name: 'Tin học',
      score1: 9,
      score2: 10,
      examScore: 9,
      totalScore: 4,
      letterScore: 'A'
    },
    {
      id: '6',
      name: 'Tiếng Anh',
      score1: 8,
      score2: 8,
      examScore: 8,
      totalScore: 2,
      letterScore: 'B'
    },
    {
      id: '7',
      name: 'Lịch sử',
      score1: 7,
      score2: 7,
      examScore: 8,
      totalScore: 3,
      letterScore: 'B'
    },
    {
      id: '8',
      name: 'Địa lý',
      score1: 8,
      score2: 8,
      examScore: 7,
      totalScore: 6,
      letterScore: 'B'
    },
    {
      id: '9',
      name: 'Giáo dục công dân',
      score1: 9,
      score2: 9,
      examScore: 8,
      totalScore: 10,
      letterScore: 'F'
    },
    {
      id: '10',
      name: 'Thể dục',
      score1: 10,
      score2: 10,
      examScore: 10,
      totalScore: 8.5,
      letterScore: 'A+'
    },
    {
      id: '11',
      name: 'Âm nhạc',
      score1: 9,
      score2: 9,
      examScore: 9,
      totalScore: 9,
      letterScore: 'D'
    }
  ]
  return (
    <div>
      <h2>Kết quả học tập</h2>
      <Separator className='my-4' />
      <ScoreView passedSubject={0} failedSubject={0} gpa={0} />
      <Separator className='my-4' />
      <TableList
        items={[
          { header: 'Tên môn học', value: 'name', className: 'min-w-[220px]' },
          { header: 'Điểm thành phần 1', value: 'score1', className: 'min-w-[80px]' },
          { header: 'Điểm thành phần 2', value: 'score2', className: 'min-w-[80px]' },
          { header: 'Điểm thi', value: 'examScore', className: 'min-w-[80px]' },
          { header: 'Điểm tổng kết', value: 'totalScore', className: 'min-w-[80px]' },
          {
            header: 'Điểm chữ',
            value: 'letterScore',
            className: 'min-w-[80px]',
            render: (item) => {
              if (item.letterScore.includes('A')) {
                return <Badge className='bg-green-500'>{item.letterScore}</Badge>
              } else if (item.letterScore.includes('B')) {
                return <Badge className='bg-blue-500'>{item.letterScore}</Badge>
              } else if (item.letterScore.includes('C')) {
                return <Badge className='bg-yellow-500'>{item.letterScore}</Badge>
              } else {
                return <Badge className='bg-red-500'>{item.letterScore}</Badge>
              }
            }
          }
        ]}
        data={data}
      />
    </div>
  )
}

export default StudentScorePage
