import { Button } from '@/components/ui/button'
import { showMessage } from '@/lib/utils/common'
import { getSignDegreeConfig } from '@/lib/utils/handle-storage'
import { KeyRound } from 'lucide-react'

const SignButton = () => {
  const signDegreeConfig = getSignDegreeConfig()
  const handleClick = () => {
    if (!signDegreeConfig.pdfSignLocation) {
      showMessage('Vui lòng cấu hình đường dẫn ứng dụng ký PDF')
      return
    }

    const url = `${signDegreeConfig.pdfSignLocation}?ts=${Date.now()}`
    window.location.href = url
  }

  return (
    <Button variant='secondary' onClick={handleClick}>
      <KeyRound />
      <span className='hidden md:block'>Ký số</span>
    </Button>
  )
}

export default SignButton
