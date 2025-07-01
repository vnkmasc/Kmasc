import React from 'react'
import Image from 'next/image'
import logoKmasc from '../../../public/assets/images/logoKMA.png'
import { Facebook, Github, GraduationCap, Instagram, MapPinned, Twitter } from 'lucide-react'
import Link from 'next/link'

// const defaultSections = [
//   {
//     title: 'Product',
//     links: [
//       { name: 'Overview', href: '#' },
//       { name: 'Pricing', href: '#' },
//       { name: 'Marketplace', href: '#' },
//       { name: 'Features', href: '#' }
//     ]
//   },
//   {
//     title: 'Company',
//     links: [
//       { name: 'About', href: '#' },
//       { name: 'Team', href: '#' },
//       { name: 'Blog', href: '#' },
//       { name: 'Careers', href: '#' }
//     ]
//   },
//   {
//     title: 'Resources',
//     links: [
//       { name: 'Help', href: '#' },
//       { name: 'Sales', href: '#' },
//       { name: 'Advertise', href: '#' },
//       { name: 'Privacy', href: '#' }
//     ]
//   }
// ]

const defaultSocialLinks = [
  { icon: <GraduationCap className='size-5' />, href: 'https://actvn.edu.vn', label: 'Trang chủ học viện' },
  { icon: <Facebook className='size-5' />, href: 'https://www.facebook.com/hocvienkythuatmatma', label: 'Facebook' },
  { icon: <MapPinned className='size-5' />, href: 'https://maps.app.goo.gl/nH4ungjtTKWfox2c8', label: 'Địa chỉ' },
  { icon: <Github className='size-5' />, href: 'https://github.com/kmasc/Kmasc', label: 'Github' }
]

const Footer: React.FC = () => {
  return (
    <footer className='w-full bg-gray-100 py-8 pt-8 dark:bg-background'>
      <div className='container space-y-4'>
        <div className='flex items-center gap-2 lg:justify-start'>
          <Link href='/'>
            <Image src={logoKmasc} alt='logo' title='Kmasc' width={32} height={32} className='h-8' />
          </Link>
          <h2 className='font-semibold text-main'>Kmasc</h2>
        </div>
        <p className='text-sm text-muted-foreground'>Giải pháp quản lý văn bằng chứng chỉ ứng dụng Blockchain.</p>
        <div className='flex flex-col justify-between gap-4 md:flex-row'>
          <p className='text-sm text-muted-foreground'>
            © 2025 Kmasc. Bản quyền thuộc về thư viện Học viện Kỹ thuật Mật mã.
          </p>
          <ul className='flex items-center space-x-6 text-muted-foreground'>
            {defaultSocialLinks.map((social, idx) => (
              <li key={idx} className='font-medium hover:text-primary'>
                <a href={social.href} aria-label={social.label} target='_blank'>
                  {social.icon}
                </a>
              </li>
            ))}
          </ul>
        </div>
      </div>
    </footer>
  )
}

export default Footer
