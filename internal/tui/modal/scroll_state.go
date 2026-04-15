package modal

// ScrollState mantém a posição do viewport em conteúdo que pode ser maior que a tela.
// É um estado mutável — pertence ao modal que o utiliza, não ao DialogFrame.
type ScrollState struct {
	// Offset é o índice da primeira linha visível no conteúdo (0-based).
	Offset int
	// Total é o número total de linhas do conteúdo.
	Total int
	// Viewport é o número de linhas visíveis (definido pelo modal em cada Render).
	Viewport int
}

// CanScrollUp retorna true se há conteúdo acima do viewport (Offset > 0).
func (s *ScrollState) CanScrollUp() bool {
	return s.Offset > 0
}

// CanScrollDown retorna true se há conteúdo abaixo do viewport.
func (s *ScrollState) CanScrollDown() bool {
	return s.Offset+s.Viewport < s.Total
}

// Up move o viewport uma linha para cima (sem ultrapassar o início).
func (s *ScrollState) Up() {
	if s.Offset > 0 {
		s.Offset--
	}
}

// Down move o viewport uma linha para baixo (sem ultrapassar o fim).
func (s *ScrollState) Down() {
	maxOffset := s.Total - s.Viewport
	if maxOffset < 0 {
		maxOffset = 0
	}
	if s.Offset < maxOffset {
		s.Offset++
	}
}

// PageUp move o viewport um viewport inteiro para cima (sem ultrapassar o início).
func (s *ScrollState) PageUp() {
	s.Offset -= s.Viewport
	if s.Offset < 0 {
		s.Offset = 0
	}
}

// PageDown move o viewport um viewport inteiro para baixo (sem ultrapassar o fim).
func (s *ScrollState) PageDown() {
	s.Offset += s.Viewport
	maxOffset := s.Total - s.Viewport
	if maxOffset < 0 {
		maxOffset = 0
	}
	if s.Offset > maxOffset {
		s.Offset = maxOffset
	}
}

// Home move o viewport para o início do conteúdo.
func (s *ScrollState) Home() {
	s.Offset = 0
}

// End move o viewport para o fim do conteúdo.
func (s *ScrollState) End() {
	maxOffset := s.Total - s.Viewport
	if maxOffset < 0 {
		maxOffset = 0
	}
	s.Offset = maxOffset
}

// ThumbLine calcula a linha (1-based dentro do viewport) onde o thumb ■ deve aparecer.
//
// Regras:
//   - Retorna -1 se o conteúdo não excede o viewport (scroll inativo).
//   - Setas têm prioridade absoluta:
//   - Se CanScrollUp() == true, a linha 1 do viewport está ocupada por ↑.
//   - Se CanScrollDown() == true, a última linha do viewport está ocupada por ↓.
//   - O thumb é posicionado proporcionalmente nas linhas restantes.
//   - Se o intervalo disponível para o thumb for zero (viewport muito pequeno), retorna -1.
func (s *ScrollState) ThumbLine() int {
	if s.Total <= s.Viewport {
		return -1
	}

	// Determinar linhas reservadas pelas setas.
	firstAvailable := 1
	lastAvailable := s.Viewport
	if s.CanScrollUp() {
		firstAvailable = 2 // linha 1 ocupada pela seta ↑
	}
	if s.CanScrollDown() {
		lastAvailable = s.Viewport - 1 // última linha ocupada pela seta ↓
	}

	available := lastAvailable - firstAvailable + 1
	if available <= 0 {
		return -1
	}

	// Posição proporcional do thumb dentro do intervalo disponível.
	// scrollable é o número máximo de passos de scroll.
	scrollable := s.Total - s.Viewport
	if scrollable == 0 {
		return firstAvailable
	}

	// Mapeia Offset → posição dentro de [0, available-1].
	thumbIndex := (s.Offset * (available - 1)) / scrollable
	return firstAvailable + thumbIndex
}
