import type React from "react"
import {Badge} from "@/components/ui/badge.tsx";
import {Wifi, Search} from "lucide-react"
import {ModeToggle} from "@/components/ThemeToggle.tsx";
import {useState} from "react";
import {Input} from "./ui/input";
import {debounce} from "@/lib/utils.ts";

interface HeaderProps {
    onSearch: (query: string) => void
    isConnected: boolean
}

export function Header({onSearch, isConnected}: HeaderProps) {
    const [searchValue, setSearchValue] = useState("")

    // Debounce search input
    const debouncedSearch = debounce((value: string) => {
        onSearch(value)
    }, 300)

    const handleSearchChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const value = e.target.value
        setSearchValue(value)
        debouncedSearch(value)
    }

    return (
        <header
            className={"h-16 px-6 flex items-center justify-between bg-[var(--background)]"}>
            <div className={"flex items-center gap-3"}>
                <h1 className={"text-3xl text-[var(--primary)]"}>
                    MQTT Dashboard
                </h1>
                {isConnected && (
                    <Badge
                        variant="outline"
                        className={"ml-2 text-green-400 dark:bg-green-300/20 dark:text-green-400 dark:border-green-400/50"}>
                        <div className={"flex items-center gap-1.5"}>
                            <Wifi className={"h-3 w-3"}/>
                            <span className={"font-medium"}>Connected</span>
                        </div>
                    </Badge>
                )}
                {
                    !isConnected && (
                        <Badge
                            variant="outline"
                            className={"ml-2 text-red-400 dark:bg-red-300/20 dark:text-red-400 dark:border-red-400/50"}>
                            <div className={"flex items-center gap-1.5"}>
                                <Wifi className={"h-3 w-3"}/>
                                <span className={"font-medium"}>Not Connected</span>
                            </div>
                        </Badge>
                    )
                }
            </div>
            <div className={"flex items-center gap-4"}>
                <div className="relative w-64">
                    <Search className="absolute left-3 top-2.5 h-4 w-4"/>
                    <Input
                        placeholder="Search topics..."
                        className="pl-9 rounded-full"
                        value={searchValue}
                        onChange={handleSearchChange}
                    />
                </div>
                <ModeToggle/>
            </div>
        </header>
    )
}