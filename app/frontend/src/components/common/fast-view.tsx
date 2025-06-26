interface Props {
  title: string
  value: any
  icon: React.ReactNode
  color?: string
}

const FastView: React.FC<Props> = (props) => {
  return (
    <div className='flex-1 rounded-lg border p-6 shadow-md'>
      <div className='flex items-center justify-between'>
        <div>
          {' '}
          <h4 className='text-gray-500'>{props.title}</h4>
          <h2 className={`font-semibold ${props.color}`}>{props.value}</h2>
        </div>
        {props.icon}
      </div>
    </div>
  )
}

export default FastView
