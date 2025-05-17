import {Button} from "../ui/button";
import {QosSelect} from "./QosSelect";
import type {QoSLevel} from "@/lib/types";

export interface TopicHeaderProps {
    topicName: string;
    isSubscribed: boolean;
    onSubscriptionToggle: () => void;
    filterQos: QoSLevel | null;
    onFilterChange: (value: string) => void;
}

export function MessageTopicHeader({
                                       topicName,
                                       isSubscribed,
                                       onSubscriptionToggle,
                                       filterQos,
                                       onFilterChange,
                                   }: TopicHeaderProps) {
    return (
        <div
            className="flex items-center justify-between p-4 bg-[var(--background)] border-y border-[var(--border)]">
            <div className="flex items-center gap-2">
                <h2>
                    {topicName}
                </h2>
                <Button
                    size="sm"
                    variant={isSubscribed ? "outline" : "default"}
                    onClick={onSubscriptionToggle}
                    className={
                        isSubscribed
                            ? "border-[#7a62f6] text-[#7a62f6] hover:bg-[#7a62f6]/10"
                            : "bg-[#7a62f6] hover:bg-[#6952e3] text-white"
                    }
                >
                    {isSubscribed ? "Unsubscribe" : "Subscribe"}
                </Button>
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
