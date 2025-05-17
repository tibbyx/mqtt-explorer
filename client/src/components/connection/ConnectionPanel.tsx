import {useState} from "react";
import {Button} from "@/components/ui/button.tsx";
import {Input} from "@/components/ui/input.tsx";
import {Label} from "@/components/ui/label.tsx";
import {Select, SelectContent, SelectItem, SelectTrigger, SelectValue} from "@/components/ui/select.tsx";

export function ConnectionPanel({onToggleConnect}: { onToggleConnect?: () => void } = {}) {
    const [host, setHost] = useState("localhost");
    const [port, setPort] = useState(1883);
    const [clientId, setClientId] = useState(`Bob-${Math.random().toString(16).substring(2, 4)}`);

    const presets = [
        {name: "Mosquitto", host: "test.mosquitto.org", port: 1883},
        {name: "HiveMQ", host: "broker.hivemq.com", port: 1883},
        {name: "EMQ X", host: "broker.emqx.io", port: 1883},
    ];

    const handlePresetSelect = (presetName: string) => {
        const preset = presets.find(p => p.name === presetName);
        if (preset) {
            setHost(preset.host);
            setPort(preset.port);
        }
    };

    const handleConnect = async () => {
        const payload = {
            Ip: host,
            Port: port.toString(),
            ClientId: clientId,
        };

        try {
            const response = await fetch("http://localhost:3000/credentials", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify(payload),
            });

            if (!response.ok) {
                throw new Error(`Connection failed with status ${response.status}`);
            }

            const result = await response.json();
            console.log("Connection successful:", result);

            if (onToggleConnect) {
                onToggleConnect();
            }

        } catch (error) {
            console.error("Connection error:", error);
            alert("Failed to connect to MQTT broker.");
        }
    };


    return (
        <div className="flex-1 flex flex-col h-full border-t overflow-auto">
            <div className="p-4 border-b flex items-center justify-between">
                <h2>
                    Connection Settings
                </h2>
            </div>

            <div className="flex flex-col justify-center items-center h-full">
                <div className="max-w-md w-full px-6 py-6 border rounded-lg">
                    <div className="space-y-4">
                        <div>
                            <Label htmlFor="preset" className="text-sm font-medium">
                                Presets
                            </Label>
                            <Select onValueChange={handlePresetSelect}>
                                <SelectTrigger id="preset" className="mt-1">
                                    <SelectValue placeholder="Select a Preset Broker"/>
                                </SelectTrigger>
                                <SelectContent>
                                    {presets.map((preset) => (
                                        <SelectItem key={preset.name} value={preset.name}>
                                            {preset.name}
                                        </SelectItem>
                                    ))}
                                </SelectContent>
                            </Select>
                        </div>

                        <div className="grid grid-cols-2 gap-4">
                            <div>
                                <Label htmlFor="host" className="text-sm font-medium">
                                    MQTT Address
                                </Label>
                                <Input
                                    id="host"
                                    className="mt-1"
                                    value={host}
                                    onChange={(e) => setHost(e.target.value)}
                                />
                            </div>

                            <div>
                                <Label htmlFor="port" className="text-sm font-medium">
                                    Port
                                </Label>
                                <Input
                                    id="port"
                                    className="mt-1"
                                    type="number"
                                    value={port}
                                    onChange={(e) => setPort(parseInt(e.target.value))}
                                />
                            </div>
                        </div>

                        <div>
                            <Label htmlFor="clientId" className="text-sm font-medium">
                                Client ID
                            </Label>
                            <Input
                                id="clientId"
                                className="mt-1"
                                value={clientId}
                                onChange={(e) => setClientId(e.target.value)}
                            />
                        </div>

                        <Button
                            className="w-full mt-6 bg-[#7a62f6] hover:bg-[#6952e3] text-white rounded-full"
                            onClick={handleConnect}
                        >
                            Connect
                        </Button>
                    </div>
                </div>
            </div>
        </div>
    );
}