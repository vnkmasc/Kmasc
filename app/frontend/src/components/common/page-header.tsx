import React from 'react'

interface Props {
  title: string
  extra?: React.ReactNode[]
}

const PageHeader: React.FC<Props> = (props) => {
  return (
    <div className='mb-4 flex items-center justify-between'>
      <h2>{props.title}</h2>
      <div className='flex items-center gap-2'>{props.extra?.map((item) => item)}</div>
    </div>
  )
}

export default PageHeader
