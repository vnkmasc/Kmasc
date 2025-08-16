import { Button } from '@/components/ui/button'
import { KeyRound } from 'lucide-react'

export default function SignButton() {
  const handleClick = () => {
    const url = `sign-pdf://test?ts=${Date.now()}`
    window.location.href = url
  }

  return (
    <Button variant='secondary' onClick={handleClick}>
      <KeyRound />
      <span className='hidden md:block'>Ký số</span>
    </Button>
  )
}
