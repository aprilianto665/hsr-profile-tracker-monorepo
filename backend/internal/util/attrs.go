package util

import (
	"fmt"
	"hsr-profile-tracker/internal/model"
)

func MergeAttributes(attrs, adds []model.Attribute) []model.Attribute {
	final := make([]model.Attribute, len(attrs))
	copy(final, attrs)
	idx := make(map[string]int, len(final))
	for i, a := range final {
		idx[a.Name] = i
	}
	for _, add := range adds {
		if i, ok := idx[add.Name]; ok {
			final[i].Value += add.Value
		} else {
			idx[add.Name] = len(final)
			final = append(final, add)
		}
	}
	return final
}

func BuildRelicSummaryOut(char model.Character) []model.RelicSummary {
	relics := make([]model.RelicSummary, 0, len(char.Relics))

	fmt.Println(char.Name)

	for _, r := range char.Relics {
		main := model.AttributeSummary{
			Name:  r.MainAffix.Name,
			Icon:  NormalizeIconPath(r.MainAffix.Icon),
			Value: FormatAttributeValue(*r.MainAffix),
		}

		subs := make([]model.AttributeSummary, 0, len(r.SubAffix))
		for _, s := range r.SubAffix {
			subs = append(subs, model.AttributeSummary{
				Name:  s.Name,
				Icon:  NormalizeIconPath(s.Icon),
				Value: FormatAttributeValue(s),
			})
		}

		score := CalculateRelicScoreValue(r, char)
		rank := GetRelicRank(score)

		relics = append(relics, model.RelicSummary{
			Name:      r.Name,
			Type:      r.Type,
			Icon:      NormalizeIconPath(r.Icon),
			Rarity:    r.Rarity,
			Level:     r.Level,
			MainAffix: &main,
			SubAffix:  subs,
			Score:     score,
			Rank:      rank,
		})
	}
	return relics
}

func BuildLightConeSummaryOut(lc *model.LightCone) *model.LightConeSummary {
	if lc == nil {
		return nil
	}

	attrs := make([]model.AttributeSummary, 0, len(lc.Attributes))
	for _, a := range lc.Attributes {
		attrs = append(attrs, model.AttributeSummary{
			Name:  a.Name,
			Icon:  NormalizeIconPath(a.Icon),
			Value: FormatAttributeValue(a),
		})
	}

	return &model.LightConeSummary{
		Name:       lc.Name,
		Rarity:     lc.Rarity,
		Rank:       lc.Rank,
		Level:      lc.Level,
		Icon:       NormalizeIconPath(lc.Icon),
		Attributes: attrs,
	}
}

func BuildFinalStatsOut(attrs, adds []model.Attribute) []model.AttributeSummary {
	final := MergeAttributes(attrs, adds)
	out := make([]model.AttributeSummary, 0, len(final))
	for _, fs := range final {
		out = append(out, model.AttributeSummary{
			Name:  fs.Name,
			Icon:  NormalizeIconPath(fs.Icon),
			Value: FormatAttributeValue(fs),
		})
	}
	return out
}

func BuildRelicScoreOut(relic []model.RelicSummary) model.RelicScore {
	var totalScore float64

	for _, r := range relic {
		totalScore += r.Score
	}

	averageScore := FloorToDecimal(totalScore/6.0, 1)

	return model.RelicScore{
		Rank:         GetRelicRank(averageScore),
		TotalScore:   FloorToDecimal(totalScore, 1),
		AverageScore: averageScore,
	}
}

func BuildProfileSummaryOut(rawData model.RawData) model.ProfileSummary {
	player := NormalizePlayerAvatar(rawData.Player)

	chars := make([]model.CharacterSummary, 0, len(rawData.Characters))
	for _, c := range rawData.Characters {

		c.Path.Icon = NormalizeIconPath(c.Path.Icon)
		c.Element.Icon = NormalizeIconPath(c.Element.Icon)

		lc := BuildLightConeSummaryOut(c.LightCone)

		relics := BuildRelicSummaryOut(c)

		relicSets := NormalizeRelicSetIcons(c.RelicSets)

		finalStats := BuildFinalStatsOut(c.Attributes, c.Additions)

		relicScore := BuildRelicScoreOut(relics)

		chars = append(chars, model.CharacterSummary{
			Name:       c.Name,
			Portrait:   NormalizeIconPath(c.Portrait),
			Rarity:     c.Rarity,
			Rank:       c.Rank,
			Level:      c.Level,
			Path:       c.Path,
			Element:    c.Element,
			LightCone:  lc,
			Relics:     relics,
			RelicSets:  relicSets,
			FinalStats: finalStats,
			RelicScore: &relicScore,
		})
	}

	summary := model.ProfileSummary{
		Player:     player,
		Characters: chars,
	}

	return summary
}
