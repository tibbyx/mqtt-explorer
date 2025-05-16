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
                className={`${
                    id ? "w-max" : "rounded-lg"
                } border-gray-200 dark:border-gray-800 bg-white dark:bg-gray-950 text-gray-700 dark:text-gray-200 focus:ring-[#7a62f6]`}
            >
                <SelectValue placeholder={label || "QoS"} className="font-medium"/>
            </SelectTrigger>
            <SelectContent className="rounded-lg border-gray-200 dark:border-gray-800">
                {showAnyOption && <SelectItem value="any">Any QoS</SelectItem>}
                <SelectItem value="0">QoS 0</SelectItem>
                <SelectItem value="1">QoS 1</SelectItem>
                <SelectItem value="2">QoS 2</SelectItem>
            </SelectContent>
        </Select>
    );
}
