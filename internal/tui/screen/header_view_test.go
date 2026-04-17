package screen

import (
	"os"
	"testing"

	"github.com/useful-toys/abditum/internal/tui/design"
	"github.com/useful-toys/abditum/internal/tui/testdata"
	"github.com/useful-toys/abditum/internal/vault"
)

// ptr retorna um ponteiro para a string fornecida.
func ptr(s string) *string {
	return &s
}

// newTestVault cria um novo vault.Manager para testes com um caminho explícito.
// O vault é inicializado com conteúdo padrão e está limpo (não modificado).
func newTestVault(caminho string) *vault.Manager {
	cofre := vault.NovoCofre()
	if err := cofre.InicializarConteudoPadrao(); err != nil {
		panic("failed to initialize test vault: " + err.Error())
	}
	return vault.NewManagerForTest(cofre, caminho)
}

// newTestVaultDirty cria um novo vault.Manager marcado como modificado.
// Marca a modificação criando uma pasta temporária e descartando-a.
func newTestVaultDirty(caminho string) *vault.Manager {
	manager := newTestVault(caminho)
	// Marcar como modificado criando uma pasta e deletando-a
	// Isso garante que IsModified() retorna true
	geral := manager.Vault().PastaGeral()
	pasta, err := manager.CriarPasta(geral, "temp", 0)
	if err != nil {
		panic("failed to create temp folder: " + err.Error())
	}
	if _, err := manager.ExcluirPasta(pasta); err != nil {
		panic("failed to delete temp folder: " + err.Error())
	}
	return manager
}

// headerRenderFn adapta HeaderView.Render para testdata.RenderFn.
// O parâmetro height é ignorado (o cabeçalho é sempre 2 linhas).
func headerRenderFn(setup func(v *HeaderView)) testdata.RenderFn {
	return func(w, _ int, theme *design.Theme) string {
		v := NewHeaderView()
		setup(v)
		return v.Render(0, w, theme)
	}
}

// --- Helper function tests ---

// goldenSizesHelpers define os tamanhos usados nos testes golden dos helpers.
var goldenSizesHelpers = []string{"10x1", "9x1", "11x1", "10x1"}

// TestRenderTab_Inactive testa a renderização de uma aba inativa.
func TestRenderTab_Inactive(t *testing.T) {
	testdata.TestRenderManaged(t, "header", "tab-inactive", []string{"10x1"},
		func(w, _ int, theme *design.Theme) string {
			rendered, _ := RenderTab("Cofre", false, theme)
			return rendered
		},
	)
}

// TestRenderTab_Active testa a renderização de uma aba ativa.
func TestRenderTab_Active(t *testing.T) {
	testdata.TestRenderManaged(t, "header", "tab-active", []string{"9x1"},
		func(w, _ int, theme *design.Theme) string {
			rendered, _ := RenderTab("Cofre", true, theme)
			return rendered
		},
	)
}

// TestRenderTabConnector_Vault testa a renderização do conector da aba Cofre.
func TestRenderTabConnector_Vault(t *testing.T) {
	testdata.TestRenderManaged(t, "header", "connector-vault", []string{"11x1"},
		func(w, _ int, theme *design.Theme) string {
			rendered, _ := RenderTabConnector("Cofre", theme)
			return rendered
		},
	)
}

// TestRenderTabConnector_Models testa a renderização do conector da aba Modelos.
func TestRenderTabConnector_Models(t *testing.T) {
	testdata.TestRenderManaged(t, "header", "connector-models", []string{"10x1"},
		func(w, _ int, theme *design.Theme) string {
			rendered, _ := RenderTabConnector("Modelos", theme)
			return rendered
		},
	)
}

// TestRenderTabConnector_Config testa a renderização do conector da aba Config.
func TestRenderTabConnector_Config(t *testing.T) {
	testdata.TestRenderManaged(t, "header", "connector-config", []string{"10x1"},
		func(w, _ int, theme *design.Theme) string {
			rendered, _ := RenderTabConnector("Config", theme)
			return rendered
		},
	)
}

// --- Component tests ---

// goldenSizesComponent define o tamanho fixo usado nos testes golden do componente.
var goldenSizesComponent = []string{"80x2"}

// TestHeader_NoVault testa a renderização sem cofre aberto.
func TestHeader_NoVault(t *testing.T) {
	testdata.TestRenderManaged(t, "header", "no-vault", goldenSizesComponent,
		headerRenderFn(func(v *HeaderView) {
			// Sem configuração: vault é nil
		}),
	)
}

// TestHeader_VaultClean testa a renderização com um cofre limpo (não modificado).
func TestHeader_VaultClean(t *testing.T) {
	testdata.TestRenderManaged(t, "header", "vault-clean", goldenSizesComponent,
		headerRenderFn(func(v *HeaderView) {
			manager := newTestVault("meu_cofre.abditum")
			v.SetVault(manager)
			v.SetActiveMode(design.WorkAreaVault)
		}),
	)
}

// TestHeader_VaultDirty testa a renderização com um cofre modificado.
func TestHeader_VaultDirty(t *testing.T) {
	testdata.TestRenderManaged(t, "header", "vault-dirty", goldenSizesComponent,
		headerRenderFn(func(v *HeaderView) {
			manager := newTestVaultDirty("meu_cofre.abditum")
			v.SetVault(manager)
			v.SetActiveMode(design.WorkAreaVault)
		}),
	)
}

// TestHeader_ModeVault testa a renderização com modo WorkAreaVault ativo.
func TestHeader_ModeVault(t *testing.T) {
	testdata.TestRenderManaged(t, "header", "mode-vault", goldenSizesComponent,
		headerRenderFn(func(v *HeaderView) {
			manager := newTestVault("meu_cofre.abditum")
			v.SetVault(manager)
			v.SetActiveMode(design.WorkAreaVault)
		}),
	)
}

// TestHeader_ModeModels testa a renderização com modo WorkAreaTemplates ativo.
func TestHeader_ModeModels(t *testing.T) {
	testdata.TestRenderManaged(t, "header", "mode-models", goldenSizesComponent,
		headerRenderFn(func(v *HeaderView) {
			manager := newTestVault("meu_cofre.abditum")
			v.SetVault(manager)
			v.SetActiveMode(design.WorkAreaTemplates)
		}),
	)
}

// TestHeader_ModeConfig testa a renderização com modo WorkAreaSettings ativo.
func TestHeader_ModeConfig(t *testing.T) {
	testdata.TestRenderManaged(t, "header", "mode-config", goldenSizesComponent,
		headerRenderFn(func(v *HeaderView) {
			manager := newTestVault("meu_cofre.abditum")
			v.SetVault(manager)
			v.SetActiveMode(design.WorkAreaSettings)
		}),
	)
}

// TestHeader_VaultNameLong testa a renderização com um nome de cofre muito longo.
func TestHeader_VaultNameLong(t *testing.T) {
	testdata.TestRenderManaged(t, "header", "vault-name-long", goldenSizesComponent,
		headerRenderFn(func(v *HeaderView) {
			manager := newTestVault("muito_muito_muito_muito_muito_muito_longo_nome_de_cofre.abditum")
			v.SetVault(manager)
			v.SetActiveMode(design.WorkAreaVault)
		}),
	)
}

// TestHeader_SearchEmpty testa a renderização com busca ativa mas sem query.
func TestHeader_SearchEmpty(t *testing.T) {
	testdata.TestRenderManaged(t, "header", "search-empty", goldenSizesComponent,
		headerRenderFn(func(v *HeaderView) {
			manager := newTestVault("meu_cofre.abditum")
			v.SetVault(manager)
			v.SetActiveMode(design.WorkAreaVault)
			v.SetSearchQuery(ptr(""))
		}),
	)
}

// TestHeader_SearchWithQuery testa a renderização com busca ativa e uma query.
func TestHeader_SearchWithQuery(t *testing.T) {
	testdata.TestRenderManaged(t, "header", "search-with-query", goldenSizesComponent,
		headerRenderFn(func(v *HeaderView) {
			manager := newTestVault("meu_cofre.abditum")
			v.SetVault(manager)
			v.SetActiveMode(design.WorkAreaVault)
			v.SetSearchQuery(ptr("senha"))
		}),
	)
}

// TestHeader_SearchQueryLong testa a renderização com uma query de busca longa.
func TestHeader_SearchQueryLong(t *testing.T) {
	testdata.TestRenderManaged(t, "header", "search-query-long", goldenSizesComponent,
		headerRenderFn(func(v *HeaderView) {
			manager := newTestVault("meu_cofre.abditum")
			v.SetVault(manager)
			v.SetActiveMode(design.WorkAreaVault)
			v.SetSearchQuery(ptr("um_termo_de_busca_muito_muito_muito_longo_para_o_espaco_disponivel"))
		}),
	)
}

// TestHeader_SearchTemplatesEmpty testa busca com aba Modelos ativa e query vazia.
func TestHeader_SearchTemplatesEmpty(t *testing.T) {
	testdata.TestRenderManaged(t, "header", "search-templates-empty", goldenSizesComponent,
		headerRenderFn(func(v *HeaderView) {
			manager := newTestVault("meu_cofre.abditum")
			v.SetVault(manager)
			v.SetActiveMode(design.WorkAreaTemplates)
			v.SetSearchQuery(ptr(""))
		}),
	)
}

// TestHeader_SearchTemplatesWithQuery testa busca com aba Modelos ativa e query.
func TestHeader_SearchTemplatesWithQuery(t *testing.T) {
	testdata.TestRenderManaged(t, "header", "search-templates-query", goldenSizesComponent,
		headerRenderFn(func(v *HeaderView) {
			manager := newTestVault("meu_cofre.abditum")
			v.SetVault(manager)
			v.SetActiveMode(design.WorkAreaTemplates)
			v.SetSearchQuery(ptr("senha"))
		}),
	)
}

// TestHeader_SearchTemplatesLong testa busca com aba Modelos ativa e query longa.
func TestHeader_SearchTemplatesLong(t *testing.T) {
	testdata.TestRenderManaged(t, "header", "search-templates-long", goldenSizesComponent,
		headerRenderFn(func(v *HeaderView) {
			manager := newTestVault("meu_cofre.abditum")
			v.SetVault(manager)
			v.SetActiveMode(design.WorkAreaTemplates)
			v.SetSearchQuery(ptr("um_termo_de_busca_muito_muito_muito_longo_para_o_espaco_disponivel"))
		}),
	)
}

// TestHeader_SearchConfigEmpty testa busca com aba Config ativa e query vazia.
func TestHeader_SearchConfigEmpty(t *testing.T) {
	testdata.TestRenderManaged(t, "header", "search-config-empty", goldenSizesComponent,
		headerRenderFn(func(v *HeaderView) {
			manager := newTestVault("meu_cofre.abditum")
			v.SetVault(manager)
			v.SetActiveMode(design.WorkAreaSettings)
			v.SetSearchQuery(ptr(""))
		}),
	)
}

// TestHeader_SearchConfigWithQuery testa busca com aba Config ativa e query.
func TestHeader_SearchConfigWithQuery(t *testing.T) {
	testdata.TestRenderManaged(t, "header", "search-config-query", goldenSizesComponent,
		headerRenderFn(func(v *HeaderView) {
			manager := newTestVault("meu_cofre.abditum")
			v.SetVault(manager)
			v.SetActiveMode(design.WorkAreaSettings)
			v.SetSearchQuery(ptr("senha"))
		}),
	)
}

// TestHeader_SearchConfigLong testa busca com aba Config ativa e query longa.
func TestHeader_SearchConfigLong(t *testing.T) {
	testdata.TestRenderManaged(t, "header", "search-config-long", goldenSizesComponent,
		headerRenderFn(func(v *HeaderView) {
			manager := newTestVault("meu_cofre.abditum")
			v.SetVault(manager)
			v.SetActiveMode(design.WorkAreaSettings)
			v.SetSearchQuery(ptr("um_termo_de_busca_muito_muito_muito_longo_para_o_espaco_disponivel"))
		}),
	)
}

func testUserHomeDir() string {
	home := os.Getenv("USERPROFILE")
	if home != "" {
		return home
	}
	home = os.Getenv("HOME")
	if home != "" {
		return home
	}
	home, _ = os.UserHomeDir()
	return home
}

// TestHeader_VaultInHomeDir testa a renderização com cofre dentro do homedir.
func TestHeader_VaultInHomeDir(t *testing.T) {
	home := testUserHomeDir()
	caminho := home + string(os.PathSeparator) + "dir1" + string(os.PathSeparator) + "dir2" + string(os.PathSeparator) + "cofre.abditum"
	testdata.TestRenderManaged(t, "header", "vault-in-home", goldenSizesComponent,
		headerRenderFn(func(v *HeaderView) {
			manager := newTestVault(caminho)
			v.SetVault(manager)
			v.SetActiveMode(design.WorkAreaVault)
		}),
	)
}

// TestHeader_VaultInHomeDirLong testa com caminho bem longo dentro do homedir.
func TestHeader_VaultInHomeDirLong(t *testing.T) {
	home := testUserHomeDir()
	caminho := home + string(os.PathSeparator) + "muito" + string(os.PathSeparator) + "muito" + string(os.PathSeparator) + "muito" + string(os.PathSeparator) + "muito" + string(os.PathSeparator) + "muito" + string(os.PathSeparator) + "muito" + string(os.PathSeparator) + "muito" + string(os.PathSeparator) + "longo" + string(os.PathSeparator) + "cofre.abditum"
	testdata.TestRenderManaged(t, "header", "vault-in-home-long", goldenSizesComponent,
		headerRenderFn(func(v *HeaderView) {
			manager := newTestVault(caminho)
			v.SetVault(manager)
			v.SetActiveMode(design.WorkAreaVault)
		}),
	)
}

// TestHeader_VaultOutsideHome testa com cofre fora do homedir.
func TestHeader_VaultOutsideHome(t *testing.T) {
	testdata.TestRenderManaged(t, "header", "vault-outside-home", goldenSizesComponent,
		headerRenderFn(func(v *HeaderView) {
			manager := newTestVault("D:/outros/cofres/meu_cofre.abditum")
			v.SetVault(manager)
			v.SetActiveMode(design.WorkAreaVault)
		}),
	)
}
