import type Resources from '@/@types/resources'

declare module 'i18next' {
    interface CustomTypeOptions {
        defaultNS: 'en'
        resources: Resources
    }
}
