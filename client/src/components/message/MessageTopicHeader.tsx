import {QosSelect} from "./QosSelect";
import type {QoSLevel} from "@/lib/types";

export interface TopicHeaderProps {
    filterQos: QoSLevel | null;
    onFilterChange: (value: string) => void;
}

export function MessageTopicHeader({
                                       filterQos,
                                       onFilterChange,
                                   }: TopicHeaderProps) {
    return (
        <div
            className="flex items-center justify-between p-4 bg-[var(--background)] border-y border-[var(--border)]">
            <div className="flex items-center gap-2">
            </div>
            <QosSelect
                value={filterQos === null ? "any" : filterQos.toString()}
                onChange={onFilterChange}
                label="Filter by QoS"
                showAnyOption
            />
        </div>
    );
}
