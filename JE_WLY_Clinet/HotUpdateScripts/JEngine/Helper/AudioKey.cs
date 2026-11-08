using System.Collections.Generic;

public class AudioKey
{

    public enum Sound
    {
        BGM_Creat,
        BGM_Hall_Rest,
        BGM_Hall_Boss,
        BGM_Hall_Battle,
        BGM_Fight,
        BGM_Fight_BeiYong,
        BGM_Fight_ShiJieBOSS,
        BGM_Fight_ShuaTuPeiYue,
        BGM_Fight_ZhanDouChangJing,
        BGM_Fight_ZhanDouPeiYue,
        BGM_Fight_ZhiYePeiYue,

        Button,
        Combo,

        Hurt_1,
        Hurt_2,
        Hurt_3,

        //Boss 攻击音效
        BossAttack,

        //近战怪物攻击
        JZAttack,

        //远程战怪物攻击
        YCAttack,

        //战斗场景BGM
        BG_ZCDT_1,
        BG_ZCDT_2,
        BG_NormalMap_1,
        BG_NormalMap_2,
        BG_NormalMap_3,
        BG_NormalMap_4,
        BG_NormalMap_5,
        BG_NormalMap_6,
        BG_NormalMap_7,

        BG_ZhanDouPeiYue,


        BG_Login,
        BG_ChooseRole1,
        BG_ChooseRole2,
        BG_Loading1,
        BG_Loading2,
        BG_Loading3,
        BG_HallRest,
        BG_HallFight,

        //UI
        UI_BossComing,
        UI_ComClick,
        UI_OpenAndCloseUI,
        UI_ComReward,
        UI_ComposeEffect,
        UI_StrengthenSuc,
        UI_StrengthStarSuc,
        UI_StrengthStarFaild,
        UI_FlyEffect,
        UI_EmailOpenAndClose,
        UI_CallEffect,
        UI_CallWhenPartnerArrive,
        UI_SellAndGotCoin,
        UI_GemEquipToSlot,
        UI_RemoveResult,
        UI_RemoveCover,
        UI_RandomCard,
        UI_ShowCard,
        UI_SelectCard,
        UI_SkillUpgradeSuc,
        UI_SkillAdvancedSuc,
        UI_SkillEquipToSlot,
        UI_PetEggCall,
        UI_PetComposed,
        UI_PetComposedWhenGotFantasticPet,
        UI_PetComposedFaild,
        UI_PetGotExp,
        UI_PurchaseGoods,
        UI_OpenRightPanel,
        UI_CloseRightPanel,
        UI_RolePowerUp,
        UI_RoleUpgrade,
        UI_BattleCDWarning,
        UI_FunctionUnlock,
        UI_EmpowermentSuc,
        UI_EmpowermentFaild,
        UI_TowerUnlock,
        UI_OpenBox,

        //战斗
        Battle_Com_Walk,
        Battle_RoleDied_Male,
        Battle_RoleDied_Female,
        Battle_MonsterHit,
        Battle_MonsterDied,

    }
    public static Dictionary<Sound, string> SoundDic = new Dictionary<Sound, string>
    {
        [Sound.Button] = "UI/UI_ComClick.mp3",//"button",
        [Sound.BGM_Hall_Rest] = "Bg_HallRest.mp3",
        [Sound.BGM_Hall_Boss] = "BG_HallBoss",//"Bg_Hall_Boss.mp3",
        [Sound.BGM_Hall_Battle] = "BG_HallFight.mp3",
        [Sound.BGM_Creat] = "Bg_Creat.mp3",
        [Sound.BGM_Fight] = "BG_Fight.mp3",
        [Sound.Combo] = "hallCombo.mp3",

        [Sound.BG_ZhanDouPeiYue] = "BG_ZhanDouPeiYue.mp3",
        //[Sound.BGM_Fight_ShiJieBOSS] = "BGM_Fight_ShiJieBOSS.mp3",
        //[Sound.BGM_Fight_ShuaTuPeiYue] = "BGM_Fight_ShuaTuPeiYue.mp3",
        //[Sound.BGM_Fight_ZhanDouChangJing] = "BGM_Fight_ZhanDouChangJing.mp3",
        //[Sound.BGM_Fight_ZhanDouPeiYue] = "BGM_Fight_ZhanDouPeiYue.mp3",
        //[Sound.BGM_Fight_ZhiYePeiYue] = "BGM_Fight_ZhiYePeiYue.mp3",

        [Sound.Hurt_1] = "Skill/Zs_tongyongshouji.mp3",
        [Sound.Hurt_2] = "Skill/Fs_tongyongshouji.mp3",
        [Sound.Hurt_3] = "Skill/Js_tongyongshouji.mp3",

        [Sound.BossAttack] = "Monster/BossAttack.mp3",
        [Sound.JZAttack] = "Monster/JZAttack.mp3",
        [Sound.YCAttack] = "Monster/YCAttack.mp3",


        [Sound.BG_ZCDT_1] = "BG_ZCDT_1.mp3",
        [Sound.BG_ZCDT_2] = "BG_ZCDT_2.mp3",
        [Sound.BG_NormalMap_1] = "BG_NormalMap_1.mp3",
        [Sound.BG_NormalMap_2] = "BG_NormalMap_2.mp3",
        [Sound.BG_NormalMap_3] = "BG_NormalMap_3.mp3",
        [Sound.BG_NormalMap_4] = "BG_NormalMap_4.mp3",
        [Sound.BG_NormalMap_5] = "BG_NormalMap_5.mp3",
        [Sound.BG_NormalMap_6] = "BG_NormalMap_6.mp3",
        [Sound.BG_NormalMap_7] = "BG_NormalMap_7.mp3",


        [Sound.BG_Login] = "BG_Login.MP3",
        [Sound.BG_ChooseRole1] = "BG_ChooseRole1.mp3",
        [Sound.BG_ChooseRole2] = "BG_ChooseRole2.mp3",
        [Sound.BG_Loading1] = "BG_Loading1.mp3",
        [Sound.BG_Loading2] = "BG_Loading2.mp3",
        [Sound.BG_Loading3] = "BG_Loading3.mp3",
        [Sound.BG_HallRest] = "BG_HallRest.mp3",
        [Sound.BG_HallFight] = "BG_HallFight.mp3",

        //UI
        [Sound.UI_BossComing] = "UI/UI_BossComing.mp3",
        [Sound.UI_ComClick] = "UI/UI_ComClick.mp3",
        [Sound.UI_OpenAndCloseUI] = "UI/UI_OpenAndCloseUI.mp3",
        [Sound.UI_ComReward] = "UI/UI_ComReward.mp3",
        [Sound.UI_ComposeEffect] = "UI/UI_ComposeEffect.MP3",
        [Sound.UI_StrengthenSuc] = "UI/UI_StrengthenSuc.MP3",
        [Sound.UI_StrengthStarSuc] = "UI/UI_StrengthStarSuc.mp3",
        [Sound.UI_StrengthStarFaild] = "UI/UI_StrengthStarFaild.mp3",
        [Sound.UI_FlyEffect] = "UI/UI_FlyEffect.mp3",
        [Sound.UI_EmailOpenAndClose] = "UI/UI_EmailOpenAndClose.mp3",
        [Sound.UI_CallEffect] = "UI/UI_CallEffect.mp3",
        [Sound.UI_CallWhenPartnerArrive] = "UI/UI_CallWhenPartnerArrive.mp3",
        [Sound.UI_SellAndGotCoin] = "UI/UI_SellAndGotCoin.MP3",
        [Sound.UI_GemEquipToSlot] = "UI/UI_GemEquipToSlot.MP3",
        [Sound.UI_RemoveResult] = "UI/UI_RemoveResult.mp3",
        [Sound.UI_RemoveCover] = "UI/UI_RemoveCover.mp3",
        [Sound.UI_RandomCard] = "UI/UI_RandomCard.MP3",
        [Sound.UI_ShowCard] = "UI/UI_ShowCard.MP3",
        [Sound.UI_SelectCard] = "UI/UI_SelectCard.MP3",
        [Sound.UI_SkillUpgradeSuc] = "UI/UI_SkillUpgradeSuc.mp3",
        [Sound.UI_SkillAdvancedSuc] = "UI/UI_SkillAdvancedSuc.mp3",
        [Sound.UI_SkillEquipToSlot] = "UI/UI_SkillEquipToSlot.mp3",
        [Sound.UI_PetEggCall] = "UI/UI_PetEggCall.mp3",
        [Sound.UI_PetComposed] = "UI/UI_PetComposed.mp3",
        [Sound.UI_PetComposedWhenGotFantasticPet] = "UI/UI_PetComposedWhenGotFantasticPet.mp3",
        [Sound.UI_PetComposedFaild] = "UI/UI_PetComposedFaild.mp3",
        [Sound.UI_PetGotExp] = "UI/UI_PetGotExp.MP3",
        [Sound.UI_PurchaseGoods] = "UI/UI_PurchaseGoods.MP3",
        [Sound.UI_OpenRightPanel] = "UI/UI_OpenRightPanel.mp3",
        [Sound.UI_CloseRightPanel] = "UI/UI_CloseRightPanel.mp3",
        [Sound.UI_RolePowerUp] = "UI/UI_RolePowerUp.MP3",
        [Sound.UI_RoleUpgrade] = "UI/UI_RoleUpgrade.MP3",
        [Sound.UI_BattleCDWarning] = "UI/UI_BattleCDWarning.mp3",
        [Sound.UI_FunctionUnlock] = "UI/UI_FunctionUnlock.mp3",
        [Sound.UI_EmpowermentSuc] = "UI/UI_EmpowermentSuc.mp3",
        [Sound.UI_EmpowermentFaild] = "UI/UI_EmpowermentFaild.MP3",
        [Sound.UI_TowerUnlock] = "UI/UI_TowerUnlock.MP3",
        [Sound.UI_OpenBox] = "UI/UI_OpenBox.MP3",


        //战斗
        [Sound.Battle_Com_Walk] = "Battle/Battle_Com_Walk.mp3",
        [Sound.Battle_RoleDied_Male] = "Battle/Battle_RoleDied_Male.mp3",
        [Sound.Battle_RoleDied_Female] = "Battle/Battle_RoleDied_Female.mp3",
        [Sound.Battle_MonsterHit] = "Battle/Battle_MonsterHit.mp3",
        [Sound.Battle_MonsterDied] = "Battle/Battle_MonsterDied.mp3",
    };

}
