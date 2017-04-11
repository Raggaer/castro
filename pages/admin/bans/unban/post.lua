function post()
    if not session:isLogged() or not session:isAdmin() then
        http:redirect("/")
        return
    end

    if http.postValues.lift_namelock then
        db:execute("DELETE FROM player_namelocks WHERE player_id = ?", http.postValues.lift_namelock)
        session:setFlash("success", "Namelock has been removed.")
        -- Clear cache
        cache:delete("SELECT * FROM player_namelocks")
    elseif http.postValues.lift_account_ban then
        db:execute("DELETE FROM account_bans WHERE account_id = ?", http.postValues.lift_account_ban)
        session:setFlash("success", "Account has been unbanned.")
        -- Clear cache
        cache:delete("SELECT * FROM account_bans")
    elseif http.postValues.lift_ip_ban then
        db:execute("DELETE FROM ip_bans WHERE ip = ?", http.postValues.lift_ip_ban)
        session:setFlash("success", "IP has been unbanned.")
        -- Clear cache
        cache:delete("SELECT * FROM ip_bans")
    end

    http:redirect("/subtopic/admin/bans")
end
