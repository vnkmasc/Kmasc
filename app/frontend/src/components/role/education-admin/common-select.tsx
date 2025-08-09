import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectLabel,
  SelectTrigger,
  SelectValue
} from '@/components/ui/select'
import { OptionType } from '@/types/common'
import { Dispatch, SetStateAction } from 'react'

interface Props {
  handleSelect: Dispatch<SetStateAction<string>>
  options: OptionType[]
  selectLabel?: string
  value: string
  placeholder?: string
}

const CommonSelect: React.FC<Props> = (props) => {
  return (
    <Select value={props.value} onValueChange={props.handleSelect}>
      <SelectTrigger>
        <SelectValue placeholder={props.placeholder} />
      </SelectTrigger>
      <SelectContent>
        <SelectGroup>
          {props.selectLabel && <SelectLabel>{props.selectLabel}</SelectLabel>}
          {props.options.map((option, index) => (
            <SelectItem key={index} value={option.value}>
              {option.label}
            </SelectItem>
          ))}
        </SelectGroup>
      </SelectContent>
    </Select>
  )
}

export default CommonSelect
