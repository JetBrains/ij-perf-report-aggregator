export interface InfoSidebarVm {
  visible: boolean

  toggle(): void
}

export class InfoSidebarVmImpl implements InfoSidebarVm {
  visible = false

  toggle(): void {
    this.visible = !this.visible
  }
}