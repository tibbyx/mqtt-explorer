import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from "../ui/select";

export interface QosSelectProps {
    id?: string;
    value: string;
    onChange: (value: string) => void;
    label?: string;
    showAnyOption: boolean;
}

export function QosSelect({
                              id,
                              value,
                              onChange,
                              label,
                              showAnyOption,
                          }: QosSelectProps) {
    return (
        <Select value={value} onValueChange={onChange}>
            <SelectTrigger
                id={id}
                size="sm"
                className="bg-white dark:bg-[var(--popover)]"
            >
                <SelectValue placeholder={label || "QoS"} className="font-medium"/>
            </SelectTrigger>
            <SelectContent>
                {showAnyOption && <SelectItem value="any">Any QoS</SelectItem>}
                <SelectItem value="0">QoS 0</SelectItem>
                <SelectItem value="1" disabled>QoS 1</SelectItem>
                <SelectItem value="2" disabled>QoS 2</SelectItem>
            </SelectContent>
        </Select>
    );
}