
export interface UserSettingsI {
    toggleShortcuts: boolean;
    shortcutsExpanded: boolean;
    acceptResponderTerms: boolean;
    doubleClickSender: boolean;
    senderPanelModal: boolean;
}


export class UserSettings implements UserSettingsI {
    toggleShortcuts: boolean;
    doubleClickSender: boolean;
    acceptResponderTerms: boolean;
    shortcutsExpanded: boolean;
    senderPanelModal: boolean;

    constructor(settings: UserSettingsI) {
        this.toggleShortcuts = settings.toggleShortcuts;
        this.doubleClickSender = settings.doubleClickSender;
        this.acceptResponderTerms = settings.acceptResponderTerms;
        this.shortcutsExpanded = settings.shortcutsExpanded;
        this.senderPanelModal = settings.senderPanelModal
    }
}