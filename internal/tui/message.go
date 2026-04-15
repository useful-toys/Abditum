package tui

// TickMsg é emitido 1 vez por segundo pelo timer global do RootModel.
// Avança a animação do spinner (MsgBusy) e decrementa o TTL de mensagens temporárias.
type TickMsg struct{}

// MessageController é a interface de controle da barra de mensagens.
// Implementada por MessageLineView em screen/. Exposta pelo RootModel via MessageController().
// Usada por ações via closure: r.MessageController().SetSuccess("Cofre salvo").
type MessageController interface {
	// SetBusy exibe spinner com texto opcional de status (ex: "Salvando...").
	// Permanente até SetXxx ou Clear explícito.
	SetBusy(text string)
	// SetSuccess exibe mensagem de sucesso. Desaparece após 5 segundos.
	SetSuccess(text string)
	// SetError exibe mensagem de erro em destaque (bold). Desaparece após 5 segundos.
	SetError(text string)
	// SetWarning exibe aviso. Desaparece após 5 segundos.
	SetWarning(text string)
	// SetInfo exibe mensagem informativa. Desaparece após 5 segundos.
	SetInfo(text string)
	// SetHintField exibe dica para o campo focado. Permanente até substituição.
	SetHintField(text string)
	// SetHintUsage exibe dica de uso geral. Permanente até substituição.
	SetHintUsage(text string)
	// Clear remove a mensagem atual imediatamente.
	Clear()
}
