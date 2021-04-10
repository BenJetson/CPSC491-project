package etsy

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"

	"github.com/BenJetson/CPSC491-project/go/app"
)

type mockHTTPClient struct {
	response string
	code     int
	err      error
}

func (m *mockHTTPClient) Do(_ *http.Request) (res *http.Response, err error) {
	if m.err != nil {
		err = m.err
		return
	}

	res = &http.Response{
		Body:       ioutil.NopCloser(strings.NewReader(m.response)),
		StatusCode: m.code,
	}
	return
}

func testsUseRealAPI(_ *testing.T) bool {
	return len(os.Getenv("TESTS_USE_REAL_ETSY_API")) > 0
}

func newTestClient(t *testing.T) (*Client, *mockHTTPClient) {
	httpClient := &mockHTTPClient{}

	c := NewClient("VOID_this_is_for_testing")
	c.httpClient = httpClient

	// Provide a way to override and use the real Etsy API for these tests.
	// This will cause the returned httpClient to be useless.
	if testsUseRealAPI(t) {
		var err error
		c, err = NewClientFromEnv()
		require.NoError(t, err)
	}

	return c, httpClient
}

func TestSearch(t *testing.T) {
	c, httpc := newTestClient(t)
	ctx := context.Background()

	t.Run("NoKeywords", func(t *testing.T) {
		httpc.code = http.StatusBadRequest
		httpc.response = "Keywords are not specific."

		_, err := c.Search(ctx, app.CommerceQuery{})
		require.Error(t, err)
	})

	t.Run("ActualQuery", func(t *testing.T) {
		resBytes, err := ioutil.ReadFile("samples/search_results.json")
		require.NoError(t, err, "failed to read sample data")

		httpc.code = http.StatusOK
		httpc.response = string(resBytes)

		actual, err := c.Search(ctx, app.CommerceQuery{
			Keywords: "Valve Portal Aperture Science Laboratories",
			Limit:    null.IntFrom(3),
		})
		require.NoError(t, err)
		assert.Len(t, actual, 3)

		if testsUseRealAPI(t) {
			return
		}

		expect := []app.CommerceProduct{
			{
				ID: 863517982,
				Title: "Aperture Laboratories T-Shirt | Orange and Blue " +
					"Aperture Labs Logo Tee | Portal inspired Gamer " +
					"Geek Apparel",
				Description: "Aperture Laboratories, also known as Aperture S" +
					"cience and Aperture Science Innovators, is a fi" +
					"ctional corporation from the Portal and Portal " +
					"2 (and also Half-Life) universe. Aperture Labor" +
					"atories is also used as a trade name by Apertur" +
					"e Science for most of its products. Before the " +
					"1970s, the corporation was known as Aperture Sc" +
					"ience Innovators.\n\nDue to a rumor that rival " +
					"company Black Mesa were developing a similar fo" +
					"rm of Portal technology, Aperture creates the f" +
					"irst Genetic Lifeform and Disk Operating System" +
					" (GLaDOS), with the long-term plan of quickenin" +
					"g the creation of the first artificial intellig" +
					"ence.\n\nAt the beginning of the 21st Century, " +
					"GLaDOS becomes fully aware and traps her creato" +
					"rs in the Aperture facility during \u0026quot;b" +
					"ring your daughter to work day,\u0026quot; view" +
					"ing her captives as lab rats for her own experi" +
					"ments. GLaDOS’ main goal was to beat Black Mesa" +
					" to the creation of Portal technology — but los" +
					"es as a result of the Black Mesa Incident which" +
					" allowed an alien force to teleport into the fa" +
					"cility, with the long-term consequence being a " +
					"war between Humans and the race known as the Co" +
					"mbine.\n\n----- ➡️ T-Shirt Details -----\n\n• U" +
					"nisex sizing, printed on premium Bella+Canvas 3" +
					"001 T-Shirts\n• Incredibly soft and stretchy, r" +
					"un true to size\n• 100% combed and ring-spun Ai" +
					"rlume cotton (heather colors contain polyester)" +
					"\n• Fabric weight: 4.2 oz (142 g/m2)\n• Crew ne" +
					"ck, shoulder-to-shoulder taping, side-seamed\n\n" +
					"This t-shirt is everything you’ve dreamed of a" +
					"nd more. It feels soft and lightweight, with th" +
					"e right amount of stretch. It’s comfortable and" +
					" flattering for both men and women. Made with p" +
					"re-shrunk 100% combed ring-spun cotton with 4.2" +
					"oz and super soft 30 singles baby jersey knit. " +
					"It has a set-in cover stitched neck and double-" +
					"needle stitching on the sleeves and hem.\n\n---" +
					"-- ➡️ Caring for your T-Shirt -----\n\nWe want " +
					"your T-Shirt to last forever, so here’s some ti" +
					"ps on looking after it.\n\n- Machine wash cold," +
					" INSIDE OUT! -- gentle cycle with mild detergen" +
					"t \u0026 similar colors.\n- Do not bleach or us" +
					"e fabric softeners.\n- Tumble dry low, or hang-" +
					"dry in the shade for longest life.\n- Cool iron" +
					" inside-out if necessary. Do not iron decoratio" +
					"n.\n- Do not dry clean.\n\n----- ➡️ Shipping an" +
					"d Production info -----\n\nOur apparel is custo" +
					"m made to order, typically taking 3-5 business " +
					"days before shipping. All parcels are sent with" +
					" tracking.\n\nDelivery estimates:\n\nUnited Sta" +
					"tes: 3-6 days after production\nCanada: 4-6 day" +
					"s after production\nUK / AU / NZ: 6-10 days aft" +
					"er production\nRest of World: 1-2 weeks\n\nIf y" +
					"ou are not in the United States we may attempt " +
					"to fulfill your order as close to your home cou" +
					"ntry as possible for faster delivery, but some " +
					"items can only be printed in the United States." +
					"\n\n(Please note: these may be longer than usua" +
					"l due to Covid-19 restrictions -- please bear w" +
					"ith us. Get in touch any time with questions ab" +
					"out an order, or for custom orders!)",
				ImageURL: null.StringFrom(
					"https://i.etsystatic.com/18828138/c/906/720/495/325/il/" +
						"752945/2793380472/il_170x135.2793380472_6g7m.jpg",
				),
				Price: app.MustMakeMoneyFromComponents(24, 95),
			},
			{
				ID:    840721432,
				Title: "Portal | Aperture Laboratories Decal",
				Description: "Portal, Aperture Laboratories Decal\n\nSize is " +
					"the width of decal.\nSee size chart picture.\n\n" +
					"Decal comes with Aperture logo in blue with bl" +
					"ack or white text, or solid color for entire de" +
					"cal (logo and text will be same color).\n\nMade" +
					" with 5-7 year weatherproof rated vinyl.\nThese" +
					" are great for outdoor or indoor applications.\n" +
					"They apply easily to any smooth, clean, nonpor" +
					"ous surface.\n\n\nHandmade from Ventureloot.com" +
					"\n\nWe make original graphic clothing inspired " +
					"by the entertainment you love and enjoy.\n\nWe " +
					"create each product by hand with careful attent" +
					"ion to detail while using the highest quality m" +
					"aterials we can.\n\nWe promise to deliver amazi" +
					"ng customer service and hope to bring entertain" +
					"ment, humor, and inspiration to every customer " +
					"we craft for.\n\nAs we strive to never stop des" +
					"igning, we will continually release new designs" +
					" for all our products. Please join us in this a" +
					"mazing venture and get your loot!\n\nOur Ventur" +
					"e Loot brand of clothing is part of our company" +
					" SQRL Ventures LLC. For custom screen printing " +
					"and decals please visit sqrlventures.com. Also," +
					" please visit our other brand White Willow Desi" +
					"gns TX. etsy.com/shop/whitewillowdesignstx",
				ImageURL: null.StringFrom(
					"https://i.etsystatic.com/21057841/d/il/32dd43/2494814508" +
						"/il_170x135.2494814508_5pbx.jpg?version=0",
				),
				Price: app.MustMakeMoneyFromComponents(4, 99),
			},
			{
				ID:    777259649,
				Title: "Aperture Laboratories logo shelf display/fridge magnet",

				// nolint: misspell // from website
				Description: "ABOUT THIS ITEM\n\nThis Aperture Laboratories l" +
					"ogo shelf display/fridge magnet has so many use" +
					"s!\n\n- Use the included stand to turn it into " +
					"a beautiful shelf display. \n- Put it on your f" +
					"ridge or microwave or use it to hold notes on y" +
					"our fridge (it will hold up to 10 sheets of A4 " +
					"paper!). \n- Put it on your car as the ultimate" +
					" retro car accessory! (NB leaving this magnet o" +
					"n your car for long periods in hot weather may " +
					"damage it) \n\nWherever you decide to display i" +
					"t, this magnet will show off your retro credent" +
					"ials for all to see!\n\nEach item is 3D-printed" +
					" in high-quality, eco-friendly PLA plastic, the" +
					"n lovingly put together by hand. During assembl" +
					"y two high-strength N52 neodymium rare-earth ma" +
					"gnets are attached securely inside (one at each" +
					" end). \n\nNo need to worry about the magnets f" +
					"alling off, being lost or chipping or scratchin" +
					"g your fridge! And because the magnets are on t" +
					"he inside, it\u0026#39;s much safer than other " +
					"fridge magnets!\n\nFEATURES\n\n- Retro-inspired" +
					" 3D-printed, multi-coloured 3D magnet\n- Includ" +
					"es free stand to turn it into a shelf display!\n" +
					"- Printed in eco-friendly PLA plastic, derived" +
					" from renewable resources.\n- Each item comes w" +
					"ith two N52 neodymium rare-earth magnets embedd" +
					"ed inside\n- Each item capable of holding up to" +
					" 10xA4 sheets of paper!\n\nYou get:\n\n-  1 x A" +
					"perture Laboratories magnet (H 45mm, W 189mm, D" +
					" 6mm)\n-  1 x Magnetic kickstand (turns it into" +
					" a shelf display)\n\n* Due to the nature of 3D " +
					"printing, each part may have small deviations f" +
					"rom the one pictured above, yet each piece is c" +
					"arefully checked for quality and the absence of" +
					" major defects! Also, expect to see some slight" +
					" evidence of the printing process on each magne" +
					"t. However, this should not detract from the ov" +
					"erall retro effect!\n\nNOTES ON DELIVERY TIMES\n" +
					"\nI will endeavour to post all items within 3 " +
					"business days. Normally I will post much quicke" +
					"r.\n\nDelivery times are as follows:\n\n- To UK" +
					" = 2-3 days.\n- To Europe = 3-5 business days.\n" +
					"- To Rest of the World (incl. USA) = 5-7 busin" +
					"ess days.",
				ImageURL: null.StringFrom(
					"https://i.etsystatic.com/17002629/d/il/f69cf0/" +
						"2183201586/il_170x135.2183201586_172c.jpg?version=0",
				),
				Price: app.MustMakeMoneyFromComponents(9, 94),
			},
		}

		assert.Equal(t, expect, actual)
	})
}

func TestGetProductByID(t *testing.T) {
	c, httpc := newTestClient(t)
	ctx := context.Background()

	// These IDs were found to be valid products as of 10th April, 2021.
	// Listings may be removed or altered such that these IDs become invalid.
	validIDs := []int{
		984592382, 886785970, 872322174, 955885411, 759011210,
		758733592, 500483967, 925349846, 802438613, 682288652,
		789812090, 743196883, 696143337, 813084067, 790804133,
		682280692, 742424653, 891018913, 543147035, 544668981,
		969193034, 792903670, 912065527, 935304707, 991621867,
	}

	t.Run("NoSuchProduct", func(t *testing.T) {
		httpc.code = http.StatusNotFound
		httpc.response = "No product by the ID given."

		_, err := c.GetProductByID(ctx, 2929292222222222)

		require.Error(t, err)
		assert.True(t, errors.Is(err, app.ErrNotFound))
	})

	t.Run("JustOne", func(t *testing.T) {
		resBytes, err := ioutil.ReadFile("samples/one_product.json")
		require.NoError(t, err, "failed to read sample data")

		httpc.code = http.StatusOK
		httpc.response = string(resBytes)

		actual, err := c.GetProductByID(ctx, validIDs[3])
		require.NoError(t, err)
		require.NotNil(t, actual)

		if testsUseRealAPI(t) {
			return
		}

		expect := app.CommerceProduct{
			ID:    840721432,
			Title: "Portal | Aperture Laboratories Decal",
			Description: "Portal, Aperture Laboratories Decal\n\nSize is " +
				"the width of decal.\nSee size chart picture.\n\n" +
				"Decal comes with Aperture logo in blue with bl" +
				"ack or white text, or solid color for entire de" +
				"cal (logo and text will be same color).\n\nMade" +
				" with 5-7 year weatherproof rated vinyl.\nThese" +
				" are great for outdoor or indoor applications.\n" +
				"They apply easily to any smooth, clean, nonpor" +
				"ous surface.\n\n\nHandmade from Ventureloot.com" +
				"\n\nWe make original graphic clothing inspired " +
				"by the entertainment you love and enjoy.\n\nWe " +
				"create each product by hand with careful attent" +
				"ion to detail while using the highest quality m" +
				"aterials we can.\n\nWe promise to deliver amazi" +
				"ng customer service and hope to bring entertain" +
				"ment, humor, and inspiration to every customer " +
				"we craft for.\n\nAs we strive to never stop des" +
				"igning, we will continually release new designs" +
				" for all our products. Please join us in this a" +
				"mazing venture and get your loot!\n\nOur Ventur" +
				"e Loot brand of clothing is part of our company" +
				" SQRL Ventures LLC. For custom screen printing " +
				"and decals please visit sqrlventures.com. Also," +
				" please visit our other brand White Willow Desi" +
				"gns TX. etsy.com/shop/whitewillowdesignstx",
			ImageURL: null.StringFrom(
				"https://i.etsystatic.com/21057841/d/il/32dd43/2494814508" +
					"/il_170x135.2494814508_5pbx.jpg?version=0",
			),
			Price: app.MustMakeMoneyFromComponents(4, 99),
		}

		assert.Equal(t, expect, actual)
	})

	t.Run("FoundMultipleSomehow", func(t *testing.T) {
		resBytes, err := ioutil.ReadFile("samples/search_results.json")
		require.NoError(t, err, "failed to read sample data")

		httpc.code = http.StatusOK
		httpc.response = string(resBytes)

		_, err = c.GetProductByID(ctx, validIDs[3])
		require.Error(t, err)
	})
}
