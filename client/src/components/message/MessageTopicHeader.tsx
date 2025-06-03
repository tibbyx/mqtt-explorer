import { Button } from "../ui/button";
import { QosSelect } from "./QosSelect";
import type { QoSLevel } from "@/lib/types";
import { X } from "lucide-react";

export interface TopicHeaderProps {
    topicName: string;
    isSubscribed: boolean;
    onSubscriptionToggle: () => void;
    filterQos: QoSLevel | null;
    onFilterChange: (value: string) => void;
    onCloseTopic: () => void;
    isSplitScreen: boolean;
    onToggleSplitScreen: () => void;
}

export function MessageTopicHeader({
                                       topicName,
                                       isSubscribed,
                                       onSubscriptionToggle,
                                       filterQos,
                                       onFilterChange,
                                       onCloseTopic,
                                       isSplitScreen,
                                       onToggleSplitScreen,
                                   }: TopicHeaderProps) {
    return (
        <div className="flex items-center justify-between p-4 bg-[var(--background)] border-y border-[var(--border)]">
            <div className="flex items-center gap-2">
                <h2>{topicName}</h2>
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
            <div className="flex items-center gap-2">
                <QosSelect
                    value={filterQos === null ? "any" : filterQos.toString()}
                    onChange={onFilterChange}
                    label="Filter by QoS"
                    showAnyOption
                />


                <Button
                    size="sm"
                    variant={isSplitScreen ? "outline" : "default"}
                    onClick={onToggleSplitScreen}
                    className="mr-2"
                    title={isSplitScreen ? "Exit Split Screen" : "Split Screen"}
                >
                    {isSplitScreen ? "Exit Split" : "Split Screen"}
                </Button>

                <button
                    onClick={onCloseTopic}
                    className="text-[#7a62f6] hover:text-[#5a42d6] p-1"
                    title="Close"
                >
                    <X size={18} />
                </button>
            </div>
        </div>
    );
}
