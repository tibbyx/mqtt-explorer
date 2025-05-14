import {useCallback, useState} from "react"

type ToastVariant = "default" | "destructive"

interface Toast {
    id: string
    title?: string
    description?: string
    variant?: ToastVariant
    className?: string
}

export function useToast() {
    const [toasts, setToasts] = useState<Toast[]>([])

    const toast = useCallback(({title, description, variant = "default", className}: Omit<Toast, "id">) => {
        const id = Math.random().toString(36).substring(2, 9)
        const newToast = {id, title, description, variant, className}

        setToasts((prevToasts) => [...prevToasts, newToast])

        // Auto dismiss after 5 seconds
        setTimeout(() => {
            setToasts((prevToasts) => prevToasts.filter((toast) => toast.id !== id))
        }, 5000)

        return id
    }, [])

    const dismiss = useCallback((id: string) => {
        setToasts((prevToasts) => prevToasts.filter((toast) => toast.id !== id))
    }, [])

    return {toast, dismiss, toasts}
}
