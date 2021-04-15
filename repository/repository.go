package repository

import (
	"context"
	"http-server/definitions"
	"log"

	"cloud.google.com/go/datastore"
)

func AddMember(client *datastore.Client, member definitions.Member) (*datastore.Key, error) {
	ctx := context.Background()
	parentKey := datastore.IDKey("Team", member.TeamID, nil)
	key := datastore.IncompleteKey("Member", parentKey)
	newMember := &definitions.Member{
		EmailAddress: member.EmailAddress,
		FirstName:    member.FirstName,
		LastName:     member.LastName,
		ID:           member.ID,
		TeamID:       member.TeamID,
	}
	return client.Put(ctx, key, newMember)
}

func AddTeam(client *datastore.Client, team *definitions.Team) (*datastore.Key, error) {
	ctx := context.Background()
	key := datastore.IncompleteKey("Team", nil)

	return client.Put(ctx, key, team)

}

func GetMembers(client *datastore.Client, teamID int64) (*[]definitions.Member, error) {
	ctx := context.Background()
	var members *[]definitions.Member
	var team definitions.Team
	err := client.Get(ctx, datastore.IDKey("Team", teamID, nil), &team)
	if err != nil {
		return nil, err
	}
	log.Println(team)
	members = &team.Members

	return members, nil
}

func RemoveMember(client *datastore.Client, id int64, teamID int64) error {
	ctx := context.Background()
	parentKey := datastore.IDKey("Team", teamID, nil)
	_, err := client.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
		if err := client.Delete(ctx, datastore.IDKey("Member", id, parentKey)); err != nil {
			return err
		}
		var team definitions.Team
		if err := tx.Get(parentKey, &team); err != nil {
			return err
		}
		log.Println(team)
		for i, m := range team.Members {
			log.Println("id to delete", id)
			log.Println("loop member id", m.ID)
			if m.ID == id {
				log.Println("member mached", m)
				log.Println("team members", team.Members)
				team.Members = append(team.Members[:i], team.Members[i+1:]...)
				log.Println("team members after deletion", team.Members)
				break
			}
		}
		log.Println(team)
		_, err := tx.Put(parentKey, &team)
		return err
	})

	return err
}

func GetMember(client *datastore.Client, id int64) (*definitions.Member, error) {
	ctx := context.Background()
	var member definitions.Member
	err := client.Get(ctx, datastore.IDKey("Member", id, nil), &member)
	member.ID = id

	return &member, err
}

func GetTeam(client *datastore.Client, id int64) (*definitions.Team, error) {
	ctx := context.Background()
	var team definitions.Team
	err := client.Get(ctx, datastore.IDKey("Team", id, nil), &team)
	team.TeamID = id

	return &team, err
}

func AddMemberToTeam(client *datastore.Client, member definitions.Member, teamID int64) (*definitions.Team, error) {
	ctx := context.Background()
	var team definitions.Team
	log.Printf("team id: %v", teamID)
	key := datastore.IDKey("Team", teamID, nil)
	_, err := client.RunInTransaction(ctx, func(tx *datastore.Transaction) error {

		if err := tx.Get(key, &team); err != nil {
			log.Println(err)
			return err
		}
		log.Println(team)
		team.Members = append(team.Members, member)
		_, err := tx.Put(key, &team)
		log.Println(team)
		return err
	})

	return &team, err
}

func GetMemberByEmail(client *datastore.Client, teamID int64, emailAddress string) (*definitions.Member, error) {
	ctx := context.Background()
	var team definitions.Team
	parentKey := datastore.IDKey("Team", teamID, nil)
	err := client.Get(ctx, parentKey, &team)
	for _, member := range team.Members {
		if member.EmailAddress == emailAddress {
			return &member, err
		}
	}
	return nil, err
}

func GetAllTeams(client *datastore.Client) ([]*definitions.Team, error) {
	ctx := context.Background()
	var teams []*definitions.Team
	query := datastore.NewQuery("Team")
	keys, err := client.GetAll(ctx, query, &teams)
	if err != nil {
		return nil, err
	}

	for i, key := range keys {
		teams[i].TeamID = key.ID
	}
	return teams, nil
}
