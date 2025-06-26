'use client'
import { type CustomFormItem } from '@/types/common'
import { Input } from '../ui/input'
import { FormDescription, FormMessage, FormControl, FormField, FormLabel, FormItem } from '../ui/form'
import { Select, SelectContent, SelectGroup, SelectItem, SelectLabel, SelectTrigger, SelectValue } from '../ui/select'
import QuerySelect from './query-select'
import { Check, ChevronsUpDown, Eye, EyeOff } from 'lucide-react'
import { Command, CommandEmpty, CommandGroup, CommandInput, CommandItem, CommandList } from '../ui/command'
import { Popover, PopoverContent, PopoverTrigger } from '../ui/popover'
import { Button } from '../ui/button'
import { cn } from '@/lib/utils'
import { useState } from 'react'

const CustomFormItem: React.FC<CustomFormItem> = (props) => {
  const [showPassword, setShowPassword] = useState(false)

  switch (props.type) {
    case 'input':
      return (
        <FormField
          control={props.control}
          name={props.name}
          render={({ field }) => (
            <FormItem>
              <FormLabel>{props.label}</FormLabel>
              <FormControl>
                <div className='relative'>
                  <Input
                    placeholder={props.placeholder}
                    {...field}
                    disabled={props.disabled}
                    type={
                      props.setting?.input?.type === 'password'
                        ? showPassword
                          ? 'text'
                          : 'password'
                        : props.setting?.input?.type || 'text'
                    }
                    onChange={(e) => {
                      if (props.setting?.input?.type === 'number') {
                        const value = e.target.value === '' ? '' : Number(e.target.value)
                        field.onChange(value)
                      } else {
                        field.onChange(e.target.value)
                      }
                    }}
                    className={props.setting?.input?.type === 'password' ? 'pr-10' : ''}
                  />
                  {props.setting?.input?.type === 'password' && (
                    <span
                      className='absolute right-3 top-1/2 -translate-y-1/2 cursor-pointer text-gray-500 hover:text-gray-700'
                      onClick={() => setShowPassword(!showPassword)}
                    >
                      {showPassword ? <EyeOff className='h-4 w-4' /> : <Eye className='h-4 w-4' />}
                    </span>
                  )}
                </div>
              </FormControl>
              {props.description && <FormDescription>{props.description}</FormDescription>}
              <FormMessage className={`${props.description ? '!mt-0' : '!mt-2'}`} />
            </FormItem>
          )}
        />
      )
    case 'select':
      return (
        <FormField
          control={props.control}
          name={props.name}
          render={({ field }) => (
            <FormItem>
              <FormLabel>{props.label}</FormLabel>
              <Select
                disabled={props.disabled}
                onValueChange={field.onChange}
                value={field.value}
                defaultValue={field.value}
              >
                <FormControl>
                  <SelectTrigger>
                    <SelectValue placeholder={props.placeholder} />
                  </SelectTrigger>
                </FormControl>
                <SelectContent>
                  {props.setting?.select?.groups?.length && props.setting?.select?.groups?.length > 0 ? (
                    props.setting?.select?.groups?.map((group, index) => (
                      <SelectGroup key={index}>
                        {group.label && <SelectLabel>{group.label}</SelectLabel>}
                        {group.options.map((option, index) => (
                          <SelectItem key={index} value={option.value}>
                            {option.label}
                          </SelectItem>
                        ))}
                      </SelectGroup>
                    ))
                  ) : (
                    <div className='flex h-16 items-center justify-center px-5 text-sm'>Không có dữ liệu</div>
                  )}
                </SelectContent>
              </Select>
              {props.description && <FormDescription>{props.description}</FormDescription>}
              <FormMessage className={`${props.description ? '!mt-0' : '!mt-2'}`} />
            </FormItem>
          )}
        />
      )
    case 'search_select':
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
                        ? props.setting?.select?.groups
                            ?.flatMap((group) => group.options)
                            ?.find((option) => option.value === field.value)?.label
                        : props.placeholder || 'Tìm kiếm và chọn'}
                      <ChevronsUpDown className='opacity-50' />
                    </Button>
                  </FormControl>
                </PopoverTrigger>
                <PopoverContent className='w-[250px] p-0'>
                  <Command>
                    <CommandInput placeholder={'Nhập để tìm kiếm'} className='h-9' />
                    <CommandList>
                      <CommandEmpty>Không tìm thấy kết quả</CommandEmpty>
                      {props.setting?.select?.groups?.map((group, groupIndex) => (
                        <CommandGroup key={groupIndex} heading={group.label}>
                          {group.options.map((option) => (
                            <CommandItem
                              key={option.value}
                              value={option.label}
                              onSelect={(currentValue) => {
                                const selectedOption = group.options.find(
                                  (opt) => opt.label.toLowerCase() === currentValue.toLowerCase()
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
    case 'query_select':
      return <QuerySelect {...props} />
    // case 'date':
    //   return (
    //     <FormField
    //       control={props.control}
    //       name={props.name}
    //       render={({ field }) => (
    //         <FormItem>
    //           <FormLabel>{props.label}</FormLabel>
    //           <Popover>
    //             <PopoverTrigger asChild>
    //               <FormControl>
    //                 <Button
    //                   variant={'outline'}
    //                   className={cn(
    //                     'w-full pl-3 text-left font-normal hover:bg-background',
    //                     !field.value && 'text-muted-foreground hover:text-muted-foreground'
    //                   )}
    //                 >
    //                   {field.value ? (
    //                     props.setting?.date?.mode === 'range' ? (
    //                       field.value.to ? (
    //                         <>
    //                           {dayjs(field.value.from).format('DD/MM/YYYY')} -{' '}
    //                           {dayjs(field.value.to).format('DD/MM/YYYY')}
    //                         </>
    //                       ) : (
    //                         dayjs(field.value).format('DD/MM/YYYY')
    //                       )
    //                     ) : (
    //                       dayjs(field.value).format('DD/MM/YYYY')
    //                     )
    //                   ) : (
    //                     <span>{props.placeholder}</span>
    //                   )}
    //                   <CalendarIcon className='ml-auto h-4 w-4 opacity-50' />
    //                 </Button>
    //               </FormControl>
    //             </PopoverTrigger>
    //             <PopoverContent className='w-auto p-0' align='start'>
    //               <Calendar
    //                 mode={props.setting?.date?.mode || 'single'}
    //                 selected={field.value}
    //                 onSelect={field.onChange}
    //                 initialFocus
    //                 disabled={(date) => {
    //                   if (props.setting?.date?.min && dayjs(date).isBefore(props.setting?.date?.min)) {
    //                     return true
    //                   }
    //                   if (props.setting?.date?.max && dayjs(date).isAfter(props.setting?.date?.max)) {
    //                     return true
    //                   }
    //                   return false
    //                 }}
    //                 numberOfMonths={props.setting?.date?.mode === 'range' ? 2 : 1}
    //               />
    //             </PopoverContent>
    //           </Popover>
    //           {props.description && <FormDescription>{props.description}</FormDescription>}
    //           <FormMessage className={`${props.description ? '!mt-0' : '!mt-2'}`} />
    //         </FormItem>
    //       )}
    //     />
    //   )
    default:
      return null
  }
}

export default CustomFormItem
