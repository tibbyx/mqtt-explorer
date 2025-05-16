import React from "react";
import {Send} from "lucide-react";
import {Textarea} from "../ui/textarea";
import {Button} from "../ui/button";
import {QosSelect} from "./QosSelect";
import type {QoSLevel} from "@/lib/types";

export interface MessageComposerProps {
    messageText: string;
    onMessageChange: (text: string) => void;
    qosLevel: QoSLevel;
    onQosChange: (qos: QoSLevel) => void;
    onPublish: (e: React.FormEvent) => void;
}

export function MessageComposer({
                                    messageText,
                                    onMessageChange,
                                    qosLevel,
                                    onQosChange,
                                    onPublish,
                                }: MessageComposerProps) {
    return (
        <div className="p-4 py-10 border-t border-gray-200 dark:border-gray-800 bg-white dark:bg-gray-950">
            <form onSubmit={onPublish} className="space-y-3">
                <div className="flex gap-4">
                    <div className="flex-1">
                        <Textarea
                            placeholder="Enter message payload..."
                            value={messageText}
                            onChange={(e) => onMessageChange(e.target.value)}
                            className="min-h-[90px] rounded-lg border-gray-200 dark:border-gray-800 focus-visible:ring-[#7a62f6] text-gray-800 dark:text-gray-200"
                        />
                    </div>
                    <div className="flex flex-col">
                        <div className="space-y-4">
                            <div className="space-y-2">
                                <QosSelect
                                    id="qos-select"
                                    value={qosLevel.toString()}
                                    onChange={(v) => onQosChange(parseInt(v) as QoSLevel)}
                                    showAnyOption={false}
                                />
                            </div>
                        </div>
                        <div className="flex justify-end">
                            <Button
                                type="submit"
                                disabled={!messageText.trim()}
                                className="w-22 h-10 bg-[#7a62f6] hover:bg-[#6952e3] text-white"
                            >
                                <Send className="h-4 w-4"/>
                                Publish
                            </Button>
                        </div>
                    </div>
                </div>

            </form>
        </div>
    );
}
