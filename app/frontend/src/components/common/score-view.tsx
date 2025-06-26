import { ChartAreaIcon, CheckCheck, OctagonAlertIcon } from 'lucide-react'
import FastView from './fast-view'
import { Badge } from '../ui/badge'

interface Props {
  passedSubject: number
  failedSubject: number
  gpa: number
  studentName?: string
  studentCode?: string
}

const ScoreView: React.FC<Props> = (props) => {
  return (
    <div className='mt-4'>
      {props.studentName && props.studentCode && (
        <div className='flex items-center gap-4'>
          <h2>{props.studentName}</h2>
          <Badge>{props.studentCode}</Badge>
        </div>
      )}
      <div className='mt-4 flex flex-col gap-4 sm:flex-row'>
        <FastView
          title='Số môn đạt'
          value={props.passedSubject ?? 0}
          icon={<CheckCheck className='text-green-500' />}
          color='text-green-500'
        />
        <FastView
          title='Số môn thi lại'
          value={props.failedSubject ?? 0}
          icon={<OctagonAlertIcon className='text-red-500' />}
          color='text-red-500'
        />
        <FastView
          title='Điểm GPA'
          value={props.gpa ?? 0}
          icon={<ChartAreaIcon className='text-blue-500' />}
          color='text-blue-500'
        />
      </div>
    </div>
  )
}

export default ScoreView
