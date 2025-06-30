import React from 'react'
import Image from 'next/image'
import { FaFacebook, FaInstagram, FaLinkedin, FaTwitter } from 'react-icons/fa'
import logoKmasc from '../../../public/assets/images/logoKMA.png'

interface Footer7Props {
  logo?: {
    url: string
    src: string
    alt: string
    title: string
  }
  sections?: Array<{
    title: string
    links: Array<{ name: string; href: string }>
  }>
  description?: string
  socialLinks?: Array<{
    icon: React.ReactElement
    href: string
    label: string
  }>
  copyright?: string
  legalLinks?: Array<{
    name: string
    href: string
  }>
}

const defaultSections = [
  {
    title: 'Product',
    links: [
      { name: 'Overview', href: '#' },
      { name: 'Pricing', href: '#' },
      { name: 'Marketplace', href: '#' },
      { name: 'Features', href: '#' }
    ]
  },
  {
    title: 'Company',
    links: [
      { name: 'About', href: '#' },
      { name: 'Team', href: '#' },
      { name: 'Blog', href: '#' },
      { name: 'Careers', href: '#' }
    ]
  },
  {
    title: 'Resources',
    links: [
      { name: 'Help', href: '#' },
      { name: 'Sales', href: '#' },
      { name: 'Advertise', href: '#' },
      { name: 'Privacy', href: '#' }
    ]
  }
]

const defaultSocialLinks = [
  { icon: <FaInstagram className='size-5' />, href: '#', label: 'Instagram' },
  { icon: <FaFacebook className='size-5' />, href: '#', label: 'Facebook' },
  { icon: <FaTwitter className='size-5' />, href: '#', label: 'Twitter' },
  { icon: <FaLinkedin className='size-5' />, href: '#', label: 'LinkedIn' }
]

const defaultLegalLinks = [
  { name: 'Terms and Conditions', href: '#' },
  { name: 'Privacy Policy', href: '#' }
]

const Footer = ({
  logo = {
    url: '/',
    src: logoKmasc.src,
    alt: 'logo',
    title: 'VnKmasc'
  },
  sections = defaultSections,
  description = 'Giải pháp quản lý văn bằng chứng chỉ ứng dụng Blockchain.',
  socialLinks = defaultSocialLinks,
  copyright = '© 2025 VnKmasc. Bản quyền thuộc về thư viện Học viện Kỹ thuật Mật mã.',
  legalLinks = defaultLegalLinks
}: Footer7Props) => {
  return (
    <footer className='mt-8 w-full bg-gray-100 pt-8 dark:bg-background'>
      <div className='container'>
        <div className='flex w-full flex-col justify-between gap-10 lg:flex-row lg:items-start lg:text-left'>
          <div className='flex w-full flex-col justify-between gap-5 lg:items-start'>
            {/* Logo */}
            <div className='flex items-center gap-2 lg:justify-start'>
              <a href={logo.url}>
                <Image src={logo.src} alt={logo.alt} title={logo.title} width={32} height={32} className='h-8' />
              </a>
              <h2 className='font-semibold text-main'>{logo.title}</h2>
            </div>
            <p className='max-w-[70%] text-sm text-muted-foreground'>{description}</p>
            <ul className='flex items-center space-x-6 text-muted-foreground'>
              {socialLinks.map((social, idx) => (
                <li key={idx} className='font-medium hover:text-primary'>
                  <a href={social.href} aria-label={social.label}>
                    {social.icon}
                  </a>
                </li>
              ))}
            </ul>
          </div>
          <div className='grid w-full gap-6 md:grid-cols-3 lg:gap-20'>
            {sections.map((section, sectionIdx) => (
              <div key={sectionIdx}>
                <h3 className='mb-4'>{section.title}</h3>
                <ul className='space-y-3 text-sm text-muted-foreground'>
                  {section.links.map((link, linkIdx) => (
                    <li key={linkIdx} className='font-medium hover:text-primary'>
                      <a href={link.href}>{link.name}</a>
                    </li>
                  ))}
                </ul>
              </div>
            ))}
          </div>
        </div>
        <div className='mt-8 flex flex-col justify-between gap-4 border-t py-8 text-xs font-medium text-muted-foreground md:flex-row md:items-center md:text-left'>
          <p className='order-2 lg:order-1'>{copyright}</p>
          <ul className='order-1 flex flex-col gap-2 md:order-2 md:flex-row'>
            {legalLinks.map((link, idx) => (
              <li key={idx} className='hover:text-primary'>
                <a href={link.href}> {link.name}</a>
              </li>
            ))}
          </ul>
        </div>
      </div>
    </footer>
  )
}

export default Footer
