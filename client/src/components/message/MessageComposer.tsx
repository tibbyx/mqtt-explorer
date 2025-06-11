import React from "react";
import {Send, X, AlertCircle} from "lucide-react";
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
    isPublishing?: boolean;
    publishError?: Error | null;
    onClearError?: () => void;
}

export function MessageComposer({
                                    messageText,
                                    onMessageChange,
                                    qosLevel,
                                    onQosChange,
                                    onPublish,
                                    isPublishing = false,
                                    publishError,
                                    onClearError,
                                }: MessageComposerProps) {
    return (
        <div className="p-4 py-10 bg-[var(--background)] border-y border-[var(--border)]">
            {publishError && (
                <div
                    className="mb-3 p-3 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-md flex items-center justify-between">
                    <div className="flex items-center gap-2 text-red-800 dark:text-red-200">
                        <AlertCircle className="h-4 w-4"/>
                        <span className="text-sm">{publishError.message}</span>
                    </div>
                    {onClearError && (
                        <Button
                            variant="ghost"
                            size="sm"
                            onClick={onClearError}
                            className="h-6 w-6 p-0 text-red-800 dark:text-red-200 hover:bg-red-100 dark:hover:bg-red-800/20"
                        >
                            <X className="h-3 w-3"/>
                        </Button>
                    )}
                </div>
            )}
            <form onSubmit={onPublish} className="space-y-3">
                <div className="flex gap-4">
                    <div className="flex-1">
                        <Textarea
                            placeholder="Enter message payload..."
                            value={messageText}
                            onChange={(e) => onMessageChange(e.target.value)}
                            className="min-h-[90px]"
                            disabled={isPublishing}
                            onKeyDown={(e) => {
                                if (e.key === "Enter" && !e.shiftKey) {
                                    e.preventDefault();
                                    const form = e.currentTarget.form;
                                    if (form) {
                                        form.requestSubmit();
                                    }
                                }
                            }}
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
                                disabled={!messageText.trim() || isPublishing}
                                className="w-24 h-10"
                            >
                                {isPublishing ? (
                                    <>
                                        <div
                                            className="animate-spin rounded-full h-4 w-4 border-b-2 border-white"></div>
                                        Sending...
                                    </>
                                ) : (
                                    <>
                                        <Send className="h-4 w-4"/>
                                        Publish
                                    </>
                                )}
                            </Button>
                        </div>
                    </div>
                </div>
            </form>
        </div>
    );
}