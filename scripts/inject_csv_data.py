import csv
import json
import os
import shutil

def main():
    base_dir = r"c:\Users\Aryan W\Desktop\postman-cli"
    users_csv = os.path.join(base_dir, "users_export.csv")
    clinics_csv = os.path.join(base_dir, "clinic_zones_export.csv")
    
    orig_env_json = os.path.join(base_dir, "vuc-prod-env.json")
    new_env_json = os.path.join(base_dir, "vuc-dynamic-env.json")
    recep_csv = os.path.join(base_dir, "receptionists_only.csv")

    # 1. Filter Receptionists into a new CSV for reqx's native --personas
    with open(users_csv, 'r', encoding='utf-8') as f_in, \
         open(recep_csv, 'w', encoding='utf-8', newline='') as f_out:
        reader = csv.DictReader(f_in)
        fieldnames = reader.fieldnames
        if not fieldnames:
            print("❌ users_export.csv has no header!")
            return
        writer = csv.DictWriter(f_out, fieldnames=fieldnames)
        writer.writeheader()
        count = 0
        for row in reader:
            if row.get('roleName') == 'RECEPTIONIST':
                writer.writerow(row)
                count += 1
    print(f"Created {recep_csv} with {count} receptionists for native --personas flag.")

    # 2. Parse Clinics
    clinics = []
    with open(clinics_csv, 'r', encoding='utf-8') as f:
        reader = csv.DictReader(f)
        for row in reader:
            clinic_id = row.get('clinicId')
            zone_ids = row.get('zoneIds', '').replace('"', '').split(',')
            if clinic_id and len(zone_ids) > 0:
                clinics.append({
                    "clinicId": clinic_id.strip(),
                    "zoneId": zone_ids[0].strip()
                })
    print(f"Loaded {len(clinics)} clinic-zone combos.")

    # 3. Create NEW Env JSON with embedded clinic data
    with open(orig_env_json, 'r', encoding='utf-8') as f:
        env_data = json.load(f)

    if 'variables' not in env_data:
        env_data['variables'] = {}
    
    env_data['variables']['clinic_data'] = json.dumps(clinics)
    env_data['name'] = "VUC Dynamic Env"

    with open(new_env_json, 'w', encoding='utf-8') as f:
        json.dump(env_data, f, indent=4)
        
    print(f"Successfully created new environment {new_env_json}")

if __name__ == "__main__":
    main()
