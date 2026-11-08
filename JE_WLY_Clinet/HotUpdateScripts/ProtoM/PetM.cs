using JEngine.Helper;
using JEngine.Singleton;
using pb;
using System.Collections.Generic;

namespace HotUpdateScripts.ProtoM
{
    /// <summary>
    /// 宠物信息
    /// </summary>
    public class PetM : SingletonTemplate<PetM>
    {
        //可获取的宠物数量
        public int[] CanGetNums = { 0, 0, 0 };
        //宠物信息
        private List<PetSlotInfo> Pets = new List<PetSlotInfo>();

        //宠物是否获取过
        private HashSet<int> PetGetHashSet = new HashSet<int>();

        //宠物道具增加时间
        public long PetRecoveryTime { get; set; }

        /// <summary>
        /// 设置宠物信息
        /// </summary>
        /// <param name="_partners"></param>
        public void SetPets(List<PetSlotInfo> _Pets)
        {
            Pets = _Pets;
        }

        public void SetGetedPets(int[] PetEverHave)
        {
            PetGetHashSet.Clear();
            if (PetEverHave != null)
            {
                for (int i = 0; i < PetEverHave.Length; i++)
                {
                    PetGetHashSet.Add(PetEverHave[i]);
                }
            }
        }

        public bool isHavePet(int cfgId)
        {
            return PetGetHashSet.Contains(cfgId);
        }

        public List<PetSlotInfo> GetPets()
        {
            return Pets;
        }

        public PetSlotInfo GetOnePetByIdx(int idx)
        {
            return Pets[idx];
        }

        public PetSlotInfo GetOnePetByCfgId(int cfgId)
        {
            PetSlotInfo petSlotInfo = null;
            for (int i = 0; i < Pets.Count; i++)
            {
                if (cfgId == Pets[i].CfgId && (petSlotInfo == null || petSlotInfo.Level < Pets[i].CfgId))
                {
                    petSlotInfo = Pets[i];
                }
            }
            return petSlotInfo;
        }

        public void SetOnePet(PetSlotInfo _partner)
        {
            bool isAdd = true;
            for (int i = 0; i < Pets.Count; i++)
            {
                if (Pets[i].Uid == _partner.Uid)
                {
                    isAdd = false;
                    Pets[i] = _partner;
                    break;
                }
            }
            if (isAdd)
                Pets.Add(_partner);
        }

        //有宠物还没上阵红点
        public int PetUPRedDot()
        {
            if (GameUtil.Instance.CheckOpenServerIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Pet, false) || GameUtil.Instance.CheakIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Pet, false))
                return 0;

            if (GetOnePetByIdx(0).CfgId != 0)
                return 0;
            for (int i = 1; i < Pets.Count; i++)
            {
                if (Pets[i].CfgId != 0)
                {
                    return 1;
                }
            }
            return 0;
        }

    }

}

