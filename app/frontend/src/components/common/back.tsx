'use client'

import { ArrowLeft } from 'lucide-react'
import { useRouter } from 'next/navigation'

const Back = () => {
  const router = useRouter()
  return <ArrowLeft className='!h-5 !w-5 cursor-pointer' onClick={() => router.back()} />
}

export default Back
