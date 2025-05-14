import type React from "react"
import {Badge} from "@/components/ui/badge.tsx";
import {Wifi, Search} from "lucide-react"
import {ModeToggle} from "@/components/ThemeToggle.tsx";
import {useState} from "react";
import {Input} from "./ui/input";
import {debounce} from "@/lib/utils.ts";

interface HeaderProps {
    onSearch: (query: string) => void
}

export function Header({onSearch}: HeaderProps) {
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
            className={"border-b border-gray-200 h-16 px-6 flex items-center justify-between bg-white shadow-sm dark:border-gray-800 dark:bg-gray-950"}>
            <div className={"flex items-center gap-3"}>
                <h1 className={"text-xl font-bold bg-gradient-to-r from-[#7a62f6] to-[#9d8bfa] bg-clip-text text-transparent"}>
                    MQTT Dashboard
                </h1>
                <Badge
                    variant="outline"
                    className={"ml-2 bg-green-50 text-green-700 dark:bg-green-900/20 dark:text-green-400 dark:border-green-800/50"}>
                    <div className={"flex items-center gap-1.5"}>
                        <Wifi className={"h-3 w-3"}/>
                        <span className={"font-medium"}>Connected</span>
                    </div>
                </Badge>

            </div>
            <div className={"flex items-center gap-4"}>
                <div className="relative w-64">
                    <Search className="absolute left-3 top-2.5 h-4 w-4 text-gray-400"/>
                    <Input
                        placeholder="Search topics..."
                        className="pl-9 border-gray-200 dark:border-gray-800 rounded-full bg-gray-50 dark:bg-gray-900 focus-visible:ring-[#7a62f6]"
                        value={searchValue}
                        onChange={handleSearchChange}
                    />
                </div>
                <ModeToggle/>
            </div>
        </header>
    )
}