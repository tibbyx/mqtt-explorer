import {MessageSquare} from "lucide-react";

export function NoTopicSelectedView() {
    return (
        <div className="flex-1 flex items-center justify-center">
            <div className="text-center p-8 max-w-md">
                <MessageSquare className="h-16 w-16 mx-auto mb-4"/>
                <h3 className="text-xl mb-2">
                    No topic selected
                </h3>
                <p>
                    Select a topic from the list to view messages or create a new topic
                    to get started.
                </p>
            </div>
        </div>
    );
}
