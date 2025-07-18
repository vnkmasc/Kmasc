import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '../ui/card'

interface ViewItemProps {
  icon: React.ReactNode
  title: string
  value: string | number | React.ReactNode
}

interface Props {
  items: ViewItemProps[]
  title: any
  description?: string
  extra?: React.ReactNode
}

const ViewItem: React.FC<ViewItemProps> = (props) => {
  return (
    <div className='flex items-center space-x-2'>
      {props.icon}
      <div>
        <p className='mb-1 text-sm font-medium'>{props.title}</p>
        <div className='text-sm text-gray-500'>
          {props.value ?? <span className='italic text-gray-500'>Không có dữ liệu</span>}{' '}
        </div>
      </div>
    </div>
  )
}

const DecriptionView: React.FC<Props> = (props) => {
  return (
    <Card>
      <div className='flex flex-row items-center p-4 sm:p-6'>
        <CardHeader className='flex-1 !p-0 sm:p-0'>
          <CardTitle>{props.title}</CardTitle>
          <CardDescription>{props.description}</CardDescription>
        </CardHeader>
        {props.extra}
      </div>
      <CardContent className='grid grid-cols-1 gap-4 sm:grid-cols-2'>
        {props.items.map((item, index) => (
          <ViewItem key={index} {...item} />
        ))}
      </CardContent>
    </Card>
  )
}

export default DecriptionView
