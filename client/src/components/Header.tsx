import type React from "react"
import {Badge} from "@/components/ui/badge.tsx";
import {Search, Wifi, User, Power, Settings} from "lucide-react";
import {ModeToggle} from "@/components/ThemeToggle.tsx";
import {useState} from "react";
import {Input} from "./ui/input";
import {debounce} from "@/lib/utils.ts";
import {Button} from "@/components/ui/button.tsx";
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuLabel,
    DropdownMenuSeparator,
    DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import {Avatar, AvatarFallback, AvatarImage} from "@/components/ui/avatar";

interface HeaderProps {
    isConnected: boolean;
    onToggleConnect: () => void
    onSearch: (value: string) => void;
}

export function Header({isConnected, onToggleConnect, onSearch}: HeaderProps) {
    const [searchValue, setSearchValue] = useState("");

    // Debounce search input
    const debouncedSearch = debounce((value: string) => {
        onSearch(value);
    }, 300);

    const handleSearchChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const value = e.target.value;
        setSearchValue(value);
        debouncedSearch(value);
    };

    return (
        <header className="h-16 px-6 flex items-center justify-between bg-[var(--background)]">
            <div className="flex items-center gap-3">
                <h1 className="text-3xl text-[var(--primary)]">
                    MQTT Dashboard
                </h1>
                {isConnected && (
                    <Badge
                        variant="outline"
                        className="ml-2 text-green-400 dark:bg-green-300/20 dark:text-green-400 dark:border-green-400/50"
                    >
                        <div className="flex items-center gap-1.5">
                            <Wifi className="h-3 w-3"/>
                            <span className="font-medium">Connected</span>
                        </div>
                    </Badge>
                )}
                {!isConnected && (
                    <Badge
                        variant="outline"
                        className="ml-2 text-red-400 dark:bg-red-300/20 dark:text-red-400 dark:border-red-400/50"
                    >
                        <div className="flex items-center gap-1.5">
                            <Wifi className="h-3 w-3"/>
                            <span className="font-medium">Not Connected</span>
                        </div>
                    </Badge>
                )}
            </div>
            <div className="flex items-center gap-4">
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
                <DropdownMenu>
                    <DropdownMenuTrigger asChild>
                        <Button variant="outline" className="relative h-9 w-9 rounded-full ">
                            <Avatar className="h-8 w-8">
                                <AvatarImage src="" alt="User"/>
                                <AvatarFallback>
                                    <User className="h-4 w-4"/>
                                </AvatarFallback>
                            </Avatar>
                        </Button>
                    </DropdownMenuTrigger>
                    <DropdownMenuContent className="w-56" align="end" forceMount>
                        <DropdownMenuLabel className="font-normal">
                            <div className="flex flex-col space-y-1">
                                <p className="text-sm font-medium leading-none">
                                    MQTT Client
                                </p>
                                <p className="text-xs leading-none text-muted-foreground">
                                    {isConnected ? "Connected to broker" : "Disconnected"}
                                </p>
                            </div>
                        </DropdownMenuLabel>
                        <DropdownMenuSeparator/>
                        {isConnected && (
                            <DropdownMenuItem onClick={onToggleConnect}>
                                <Power className="mr-2 h-4 w-4"/>
                                <span>Disconnect</span>
                            </DropdownMenuItem>
                        )}
                        <DropdownMenuItem>
                            <Settings className="mr-2 h-4 w-4"/>
                            <span>Settings</span>
                        </DropdownMenuItem>
                    </DropdownMenuContent>
                </DropdownMenu>
            </div>
        </header>
    );
}