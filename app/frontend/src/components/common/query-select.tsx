'use client'

import { Check, ChevronsUpDown, Loader2 } from 'lucide-react'
import { CommandGroup, CommandItem, CommandEmpty, CommandInput, CommandList, Command } from '../ui/command'
import { PopoverContent, PopoverTrigger } from '../ui/popover'
import { FormControl, FormField, FormItem, FormLabel, FormDescription, FormMessage } from '../ui/form'
import { Popover } from '@radix-ui/react-popover'
import { Button } from '../ui/button'
import { cn } from '@/lib/utils'
import { CustomFormItem } from '@/types/common'
import { debounce } from '@/lib/utils/common'
import useSWR from 'swr'
import { useState } from 'react'
import { usePathname } from 'next/navigation'

const QuerySelect: React.FC<CustomFormItem> = (props) => {
  const [searchText, setSearchText] = useState('')
  const pathname = usePathname()
  const { data, isLoading } = useSWR(pathname + searchText, () => props.setting?.querySelect?.queryFn(searchText), {
    isPaused() {
      return !searchText
    }
  })
  return (
    <FormField
      control={props.control}
      name={props.name}
      render={({ field }) => (
        <FormItem>
          <FormLabel className=''>{props.label}</FormLabel>
          <Popover>
            <PopoverTrigger asChild>
              <FormControl>
                <Button
                  variant='outline'
                  role='combobox'
                  className={cn(
                    'w-full justify-between px-3 py-1 hover:bg-background',
                    !field.value && 'text-muted-foreground hover:text-muted-foreground'
                  )}
                  disabled={props.disabled}
                >
                  {field.value
                    ? (data ?? []).find((item: any) => item.value === field.value)?.label || field.value
                    : props.placeholder || 'Tìm kiếm và chọn'}
                  {isLoading ? <Loader2 className='animate-spin' /> : <ChevronsUpDown className='opacity-50' />}
                </Button>
              </FormControl>
            </PopoverTrigger>
            <PopoverContent className='min-w-[250px] p-0'>
              <Command>
                <CommandInput
                  placeholder={'Nhập để tìm kiếm'}
                  className='h-9'
                  onChangeCapture={debounce((e: any) => {
                    setSearchText(e.target.value)
                  }, 200)}
                />
                <CommandList>
                  <CommandEmpty>Không tìm thấy kết quả</CommandEmpty>
                  {(data ?? []).map((item: any, itemIndex: any) => (
                    <CommandGroup key={itemIndex} heading={item.label}>
                      {item.options.map((option: any) => (
                        <CommandItem
                          key={option.value}
                          value={option.label}
                          onSelect={(currentValue) => {
                            const selectedOption = item.options.find(
                              (opt: any) => opt.label.toLowerCase() === currentValue.toLowerCase()
                            )
                            if (selectedOption) {
                              field.onChange(selectedOption.value === field.value ? '' : selectedOption.value)
                            }
                          }}
                        >
                          {option.label}
                          <Check
                            className={cn('ml-auto', field.value === option.value ? 'opacity-100' : 'opacity-0')}
                          />
                        </CommandItem>
                      ))}
                    </CommandGroup>
                  ))}
                </CommandList>
              </Command>
            </PopoverContent>
          </Popover>
          {props.description && <FormDescription>{props.description}</FormDescription>}
          <FormMessage className={`${props.description ? '!mt-0' : '!mt-2'}`} />
        </FormItem>
      )}
    />
  )
}

export default QuerySelect
